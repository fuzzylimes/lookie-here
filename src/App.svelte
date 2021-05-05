<script>
import { afterUpdate } from "svelte";
import Popup from './Popup.svelte';
import Landing from './Landing.svelte';

	let display = false;
	let map_url, map, markers, link, modal;

	const urlParams = new URLSearchParams(window.location.search);
	const query = urlParams.get('q');
	if (query) display = true;
	
	afterUpdate(() => {
		if (display) {
			setmap();
		}
	}) 

	const setURL = ({detail}) => {
		map_url = detail.url
		display = true;
	}
	
	const encodeURL = '/api/map?m=encode';
	const decodeURL = '/api/map?m=decode';
	
	async function decode() {
		// console.log("Query is set");
		return fetch(decodeURL, {
			headers: { "content-type": "text/plain" },
			body: query,
			method: "POST"
		})
		.then(data => { return data.json() })
		.then(res => {
			// console.log(res);
			map_url = res.url;
			markers = res.markers;
		})
		.catch(error => { 
			console.error(error);
			document.location.href = '/';
		}) 
	}

	async function setmap() {
		if(query) {
			await decode();
		}
		configure_map()
		// console.log("image set to url")

		// Handle export button
		var ex = document.getElementById('export');
		ex.addEventListener('click', (event) => {
			event.stopPropagation();
			let marks = buildMapping();
			// console.log(marks);
			fetch(encodeURL, {
				headers: {"content-type": "application/json"},
				body: marks,
				method: "POST"
			})
			.then(data => {return data.json()})
			.then(res => {
				// console.log(res);
				link = `${window.location.host}/?q=${res.d}`;
				showModal();
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

		// var message_text = document.getElementById("modal_text");

		modal = document.querySelector('.modal');
		const close = document.querySelector('.modal-close')
	
		close.addEventListener('click',function () {
			showModal();
		})
	
		window.addEventListener('click',function (event) {
			if (event.target.className === 'modal-background') {
				showModal();
			}
		})
	}

	function configure_map() {
			// console.log("setmap triggered")
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
		// console.log(`removing ${id}`)
		let new_markers = []
		markerObs.forEach(marker => {
			if(marker._id === id) {
				map.removeLayer(marker)
			} else {
				new_markers.push(marker)
			}
		})
		markerObs = new_markers;
		// console.log(markerObs)
	}

	function addToMap(id, latlng, note) {

		let marker = L.marker(latlng);
		bindPopup(marker, (m) => {
			let c = new Popup({
				target: m,
				props: {
					note: note ? note : marker.note,
				},
				maxHeight: 225,
				maxWidth: 255,
			});
			c.$on('change', ({detail}) => {
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
		// console.log(markerObs)
	}

	function updateNote(id, note) {
		// console.log(id, note)
		markerObs.forEach(marker => {
			if(marker._id === id) {
				marker.note = note;
			}
		});
		// console.log(markerObs)
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

	function onMapClick(e) {
		let popups = document.querySelector('.leaflet-popup-close-button');
		if (popups) {
			popups.click();
			return;
		}
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
				latlong: [Math.ceil(m._latlng.lat), Math.ceil(m._latlng.lng)],
				note: m.note
			})
		})
		// console.log(marks)
		let payload = {
			url: map_url,
			markers: marks
		}
		return JSON.stringify(payload);
	}

	function ToClipboard () {
		var dummy = document.createElement("textarea");
		document.body.appendChild(dummy);
		dummy.value = link;
		dummy.select();
		document.execCommand("copy");
		document.body.removeChild(dummy);
		showModal();
		bulmaToast.toast({ 
			message: 'Copied to clipboard!',
			type: 'is-link',
  			position: 'center'
		});
	}

	function showModal() {
		document.querySelector('html').classList.toggle('is-clipped');
		modal.classList.toggle("is-active");
	}

</script>

<main>

	{#if display}
	<div class="map">
        <div id="mapid">
			<div class="map-button export">
				<button aria-label="Export Map" id="export" class="button is-link"><i class="fas fa-file-export"></i></button>
			</div>
			<div class="map-button reset">
				<button id="reset" class="button is-danger"><i class="fas fa-undo"></i></button>
			</div>
        </div>
    </div>
    <!-- The Modal -->
    <div id="myModal" class="modal has-text-centered">
		<div class="modal-background"></div>
		<div class="modal-content">
			<h1 class="title">Sharable Link</h1>
			<div class="button is-link" on:click={ToClipboard}>Copy To Clipboard</div>
		</div>
		<button class="modal-close is-large" aria-label="close"></button>
	</div>
	{:else}
	<div>
		<Landing on:change={setURL}/>
	</div>
	{/if}
	
	
</main>

<style>
	.map {
		z-index: -1;
	}

	.map-button{
		display: flex;
		align-items: center;
		position: absolute;
		opacity: 1;
		text-align: center;
		z-index: 1001;
	}
	.map-button:hover{
		opacity: .9;
		cursor: pointer;
	}
	.export {
		top: 10px;
		right: 10px;
		/* width: 50px; */
		height: 40px;
	}

	.export .button {
	max-width: 50px;
	}
	.reset {
		top: 60px;
		right: 10px;
		/* width: 70px; */
		height: 40px;
	}

	/* Modal Content */
	.modal-content {
		background-color: #fefefe;
		margin: auto;
		padding: 20px;
		border: 1px solid #888;
		border-radius: 15px;
		/* width: 50vw; */
		z-index: 1002;
	}

	.modal-background {
		height: 100vh;
		width: 100vw;
	}

	.modal {
		z-index: 1005;
		background-color: aliceblue;
	}
	
</style>