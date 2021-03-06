<script>
import { afterUpdate } from "svelte";
import Popup from './Popup.svelte';
import Landing from './Landing.svelte';

	let display = false;
	let map_url, map, markers, link, export_modal, help_modal;

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
				link = `${window.location.protocol}//${window.location.host}/?q=${res.d}`;
				toggleModal(export_modal);
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

		export_modal = document.querySelector('#export-modal');
		help_modal = document.querySelector('#help-modal');
		const close_export = document.querySelector('#export-modal-close')
		const close_help = document.querySelector('#help-modal-close')
	
		close_export.addEventListener('click',function () {
			toggleModal(export_modal);
		})

		close_help.addEventListener('click',function () {
			toggleModal(help_modal);
		})
	
		window.addEventListener('click',function (event) {
			if (event.target.className.includes('modal-background')) {
				toggleModal(event.target.parentNode);
			}
		})
	}

	function configure_map() {
			// console.log("setmap triggered")
			map = L.map('mapid', {
				crs: L.CRS.Simple,
				minZoom: -2,
				maxZoom: 2,
				zoom: -1,
				center: [0, 0]
			});
			const img = new Image();
			img.onload = function() {
				var bounds = [[0, 0], [this.height, this.width]];
				var image = L.imageOverlay(map_url, bounds).addTo(map);
				map.panTo(new L.LatLng(Math.ceil(this.height/2), Math.ceil(this.width/2)));
			}
			img.src = map_url;
			
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
		toggleModal(export_modal);
		bulmaToast.toast({ 
			message: 'Copied to clipboard!',
			type: 'is-link',
  			position: 'center'
		});
	}

	function toggleModal(e) {
		document.querySelector('html').classList.toggle('is-clipped');
		e.classList.toggle("is-active");
	}

	function toggleHelp(e) {
		e.stopPropagation();
		toggleModal(help_modal);
	}

</script>

<main>

	{#if display}
	<div class="map">
        <div id="mapid">
			<div class="map-button export">
				<button aria-label="Export Map" id="export" class="button is-link"><i class="fas fa-link"></i></button>
			</div>
			<div class="map-button reset">
				<button id="reset" class="button is-danger"><i class="fas fa-undo"></i></button>
			</div>
			<div class="map-button help">
				<button id="help" class="button is-info" on:click={toggleHelp}><i class="far fa-question-circle"></i></button>
			</div>
        </div>
    </div>
    <!-- The Modal -->
    <div id="export-modal" class="modal has-text-centered">
		<div class="modal-background"></div>
		<div class="modal-content">
			<h1 class="title">Sharable Link</h1>
			<div class="button is-link" on:click={ToClipboard}>Copy To Clipboard</div>
		</div>
		<button id="export-modal-close" class="modal-close is-large" aria-label="close"></button>
	</div>
	<div class="modal" id="help-modal">
		<div class="modal-background"></div>
		<div class="modal-card">
			<header class="modal-card-head">
				<p class="modal-card-title">Lookie-Here - Help</p>
			</header>
			<section class="modal-card-body">
				<p class="title is-3">Placing Pins</p>
				<p>To place pins, simply click/tap on the screen. A pin shoud appear under where you clicked/tapped.
				</p>
				<p class="title is-3 mt-3">Adding Notes</p>
				<p>You can add optional notes to your pins. Click/tap on a pin after it has been placed and then click/tap in the text box to add a note. Enter anything you want and press save. Notes will be included when sharing.</p>
				<p class="title is-3 mt-3">Deleting Pins</p>
				<p>Individual pins can be deleted by double clicking/tapping on a pin.</p>
				<p>All pins can be removed by clicking on the <i class="fas fa-undo"></i> button.</p>
				<p class="title is-3 mt-3">Sharing Link</p>
				<p>To generate a link to your map, click on the <i class="fas fa-link"></i> button, and press the "Copy to Clipboard" button in the pop up. Your sharable link will be in your active clipboard to be pasted.</p>
			</section>
			<footer class="modal-card-foot">
			<button id="help-modal-close" class="button">Close</button>
			</footer>
		</div>
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
		right: 10px;
		height: 40px;
	}
	.map-button:hover{
		opacity: .9;
		cursor: pointer;
	}
	.export {
		top: 10px;
	}
	.export .button {
	max-width: 50px;
	}
	.reset {
		top: 60px;
	}
	.help {
		top: 110px;
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
		height: 100%;
		width: 100%;
		z-index: 1001;
	}

	.modal {
		z-index: 1005;
	}

	.modal-card {
		z-index: inherit;
	}
	
</style>