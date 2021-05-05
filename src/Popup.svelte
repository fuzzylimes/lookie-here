<script>
    import { createEventDispatcher } from 'svelte';
	const dispatch = createEventDispatcher();

    export let note;
    let edit = false;

    const note_update = () => {
            dispatch('change', note)
            edit = ! edit;
        }
</script>

<style>
    .text-area {
        text-align: center;
    }
    .scroll {
        white-space: pre-line;
        min-width: 150px;
    }

</style>

<div>
    <div class="text-area">
        {#if edit}
        <textarea class="textarea" placeholder="Enter any details" bind:value={note}></textarea>
        <div class="button is-small is-link is-light mt-1" on:click={note_update}>Save</div>
        {:else}
        <div on:click={() => edit = true} class="scroll">
            {note ? note : "Click to add a note"}
        </div>
        {/if}
    </div>
</div>