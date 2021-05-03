<script>
    import { createEventDispatcher } from 'svelte';
	const dispatch = createEventDispatcher();

    export let note;
    let edit = false;

    const note_update = () => {
            console.log("Hit!")
            dispatch('change', note)
        }
    
    // document.addEventListener('DOMContentLoaded', function(){
    //     const text_tag = document.getElementById("txt")
    //     text_tag.addEventListener('click', () => {
    //         edit = true;
    //     })
    // });
</script>

<style>
    .text-area {
        padding: .5rem;
        text-align: center;
    }
</style>

<div>
    <div class="text-area">
        {#if edit}
        <textarea placeholder="Enter any details" bind:value={note} on:keypress={note_update}></textarea>
        {:else}
        <div on:click={() => edit = true} class="scroll">
            {note ? note : "Click to add a note"}
        </div>
        {/if}
    </div>
</div>