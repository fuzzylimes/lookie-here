<script>
    import { createEventDispatcher, onMount } from 'svelte'
    const dispatch = createEventDispatcher();

    function handleURL() {
        dispatch('change', {
            url: c ? document.getElementById('custom_url').value : selected
        })
    }

    let selected;
    $: c = selected == "custom";
</script>

<style>

    :global(body) {
        background-color: black;
        color: white;
    }

    .selection {
        font-size: 2rem;
        text-align: center;
    }

    .center {
        height: 100vh;
        width: 100vw;
        display: flex;
        align-items: center;
        justify-content: center;
    }

    .drop-down {
        font-size: 2rem;
    }


</style>

<div class="center">
    <div class="selection">
        <label for="sites">Choose a map:</label>
        <div class="drop-down">
            <select bind:value={selected} name="sites" id="sites">
                <option value="custom">Custom URL</option>
                <option value="https://i.redd.it/57ic5mnuza861.png">Genshin Impact</option>
                <option value="https://wallpaperaccess.com/full/815670.jpg">Mordor</option>
                <option value="https://i.pinimg.com/originals/b0/fa/34/b0fa349714c4aaf4fdc843ebf5df6be7.jpg">LoZ: BotW</option>
            </select>
        </div>
        {#if c}
        <div class="custom">
            <input id="custom_url" type="url" placeholder="Custom map URL">
        </div>
        {/if}
        <button type="button" on:click={handleURL}>Submit</button>
    </div>
</div>