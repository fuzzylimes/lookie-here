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
    #custom_url {
        max-width: 60%;
    }

    .header-image {
        position: relative;
        background-image: url("https://unsplash.com/photos/dWwZBcjw3GE/download?force=true&w=1920");
        background-position: center center;
        background-repeat: no-repeat;
        background-attachment: fixed;
        background-size: cover;
        background-color: #999;
    }

    .header-image::before {
    content: "";
    position: absolute;
    top: 0px;
    right: 0px;
    bottom: 0px;
    left: 0px;
    background-color: rgba(0,0,0,0.6);
}

</style>

<section class="hero is-fullheight header-image">
    <div class="hero-body">
        <div class="container">
            <div class="columns is-vcentered is-centered">
                <div class="column is-half-tablet is two-thirds-widescreen has-text-centered mb-5-mobile">
                    <h1 class="title is-1 is-spaced has-text-white is-uppercase">
                        Lookie-Here!
                    </h1>
                    <div class="subtitle is-4 has-text-white">
                        <p>
                            Share map markers with friends on any map you want!
                        </p>
                        <p>
                            Use one of our default maps or use your own!
                        </p>
                    </div>
                </div>
                <div class="column is-half-tablet is two-thirds-widescreen has-text-centered">
                    <div class="container">
                        <div class="select mb-4">
                            <select bind:value={selected} name="sites" id="sites">
                                <option value="custom">Custom URL</option>
                                <option value="https://i.redd.it/57ic5mnuza861.png">Genshin Impact</option>
                                <option value="https://wallpaperaccess.com/full/815670.jpg">Mordor</option>
                                <option value="https://i.pinimg.com/originals/b0/fa/34/b0fa349714c4aaf4fdc843ebf5df6be7.jpg">LoZ: BotW</option>
                            </select>
                        </div>
                        {#if c}
                        <input class="input" id="custom_url" type="url" placeholder="Custom map URL">
                        {/if}
                        <div class="container pt-4">
                            <button class="button is-link" type="button" on:click={handleURL}>Get Started!</button>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</section>