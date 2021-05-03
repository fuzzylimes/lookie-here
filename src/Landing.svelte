<script>
    import { createEventDispatcher } from 'svelte'
    const dispatch = createEventDispatcher();

    function handleURL() {
        let val = document.getElementById("sites").value;
        dispatch('change', {
            url: val
        })
    }
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
            <select name="sites" id="sites">
                <option value="https://i.redd.it/57ic5mnuza861.png">Genshin Impact</option>
                <option value="saab">Saab</option>
                <option value="mercedes">Mercedes</option>
                <option value="audi">Audi</option>
            </select>
        </div>
        <button type="button" on:click={handleURL}>Submit</button>
    </div>
</div>