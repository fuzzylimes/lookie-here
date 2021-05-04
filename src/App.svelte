<script>
import { afterUpdate } from "svelte";
import Popup from './Popup.svelte';
import Landing from './Landing.svelte';

	let display = false;
	let map_url, map, markers;

	const urlParams = new URLSearchParams(window.location.search);
	const query = urlParams.get('q');
	if (query) display = true;
	
	afterUpdate(() => {
		console.log("Update Triggered")
		if (display) {
			setmap();
		}
	}) 

	const setURL = ({detail}) => {
		map_url = detail.url
		display = true;
	}
	
	const encodeURL = '/api/encode';
	const decodeURL = '/api/decode';
	
	async function decode() {
		console.log("Query is set");
		return fetch(decodeURL, {
			headers: { "content-type": "text/plain" },
			body: query,
			method: "POST"
		})
		.then(data => { return data.json() })
		.then(res => {
			console.log(res);
			map_url = res.url;
			markers = res.markers;
		})
		.catch(error => { console.error(error) }) 
	}

	async function setmap() {
		if(query) {
			await decode();
		}
		configure_map()
		console.log("image set to url")

		// Handle export button
		var ex = document.getElementById('export');
		ex.addEventListener('click', (event) => {
			event.stopPropagation();
			let marks = buildMapping();
			console.log(marks);
			fetch(encodeURL, {
				headers: {"content-type": "application/json"},
				body: marks,
				method: "POST"
			})
			.then(data => {return data.json()})
			.then(res => {
				console.log(res);
				message_text.innerHTML = `<a href="/?q=${res.d}">Copy Link</a>`
				modal.style.display = "block";
			})
			.catch(error => {console.error(error)})
		});

		// Handle reset button
		var reset = document.getElementById('reset');
		reset.addEventListener('click', (event) => {
			event.stopPropagation();
			clearAllMarkers();
			window.history.replaceState({}, document.title, "/");
		});

		map.on('click', onMapClick);

		var modal = document.getElementById("myModal");

		// Get the button that opens the modal
		var btn = document.getElementById("myBtn");

		// Get the <span> element that closes the modal
		var span = document.getElementsByClassName("close")[0];

		var message_text = document.getElementById("modal_text");

		// When the user clicks on <span> (x), close the modal
		span.onclick = function () {
			modal.style.display = "none";
		}

		// When the user clicks anywhere outside of the modal, close it
		window.onclick = function (event) {
			event.stopPropagation();
			if (event.target == modal) {
				modal.style.display = "none";
			}
		}
	}

	function configure_map() {
			console.log("setmap triggered")
			map = L.map('mapid', {
				crs: L.CRS.Simple,
				minZoom: -1,
				maxZoom: 2,
				zoom: 0,
				center: [1000, 1200]
			});
			var bounds = [[0, 0], [2110, 2469]];
			// var image = L.imageOverlay("https://i.redd.it/57ic5mnuza861.png", bounds).addTo(map);
			var image = L.imageOverlay(map_url, bounds).addTo(map);

			// Add in the markers
			if(markers) {
				markers.forEach(element => {
					addToMap(element.id, element.latlong, element.note)
				});
			}
	}

	let markerObs = []

	// Handle query
	function clearMarker(id) {
		console.log(`removing ${id}`)
		let new_markers = []
		markerObs.forEach(marker => {
			if(marker._id === id) {
				map.removeLayer(marker)
			} else {
				new_markers.push(marker)
			}
		})
		markerObs = new_markers;
		console.log(markerObs)
	}

	function addToMap(id, latlng, note) {
		// let n = `${note}<br><button on:click="{() => clearMarker(${id})}">Remove</button>`;
		// let marker = L.marker(latlng).addTo(map).bindPopup(n);

		let marker = L.marker(latlng);
		bindPopup(marker, (m) => {
			let c = new Popup({
				target: m,
				props: {
					note: note ? note : marker.note,
				}
			});
			c.$on('change', ({detail}) => {
				console.log("THIS CHANGED")
				updateNote(id, detail);
			});
			return c;
		});
		marker.addTo(map);
		marker._id = id;
		marker.on('dblclick',function(e) {
			clearMarker(this._id);
		});
		markerObs.push(marker)
		console.log(markerObs)
	}

	function updateNote(id, note) {
		console.log(id, note)
		markerObs.forEach(marker => {
			if(marker._id === id) {
				marker.note = note;
			}
		});
		console.log(markerObs)
	}

	function bindPopup(marker, createFn) {
		let popupComponent;
		marker.bindPopup(() => {
			let container = L.DomUtil.create('div');
			popupComponent = createFn(container);
			return container;
		});

		marker.on('popupclose', () => {
			if(popupComponent) {
				let old = popupComponent;
				popupComponent = null;
				// Wait to destroy until after the fadeout completes.
				setTimeout(() => {
					old.$destroy();
				}, 500);

			}
		});
	}

	function clearAllMarkers() {
		markerObs.forEach(marker => {
			map.removeLayer(marker)
		})
		markerObs = []
	}

	// map.setView(, -1);

	function onMapClick(e) {
		// marker = L.marker(e.latlng).addTo(map)
		//     .bindPopup("<b>Hello world!</b><br />I am a popup.").openPopup();
		addToMap(getNextId(), e.latlng, null)
	}

	function getNextId() {
		let ids = markerObs.map(e => e._id)
		for (let count = 1; count <= markerObs.length; count++) {
		if (!ids.includes(count)) {
			return count
		}
		}
		return markerObs.length + 1;
	}

	function buildMapping() {
		let marks = []
		markerObs.forEach(m => {
			marks.push({
				id: m._id,
				latlong: [m._latlng.lat, m._latlng.lng],
				note: m.note
			})
		})
		console.log(marks)
		let payload = {
			url: map_url,
			markers: marks
		}
		return JSON.stringify(payload);
		// return btoa(JSON.stringify(marks)).replace(/\+/g, '-').replace(/\//g, '_').replace(/=+$/g, '');
	}
</script>

<main>

	{#if display}
	<div class="map">
        <div id="mapid">
			<div class="map-button export">
				<button id="export" class="button is-link"><i class="fas fa-file-export"></i></button>
			</div>
			<div class="map-button reset">
				<button id="reset" class="button is-danger"><i class="fas fa-undo"></i></button>
			</div>
        </div>
    </div>
    <!-- The Modal -->
    <div id="myModal" class="modal">
    
        <!-- Modal content -->
        <div class="modal-content">
            <span class="close">&times;</span>
            <p id="modal_text">Some text in the Modal..</p>
        </div>
    
    </div>
	{:else}
	<div>
		<Landing on:change={setURL}/>
	</div>
	{/if}
	
	
</main>

<style>
	
</style>