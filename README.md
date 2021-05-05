# Lookie-Here!

Lookie-Here provides a way to share maps with markers with other people. It was created to provide an
easier way of sharing maps for fantasy worlds between people, much like that of Google Maps.

The application allows you to provide a URL to any map you want (or use one of the included options),
and place markers on that map. You then have the ability to generate a link to that map and share with 
anyone!

## Usage
### Selecting a Map
Lookie-Here relies on publicly available map images in order to work. You have one of two options:
1. Use one of our provided maps
2. Supply a URL to your own map

> The higher quality an image, the better! Anything below 2k will still work, but not be as effective.

### Using the Map
After you have selected a map to use, the map will be displayed. If you have used any kind of map application
before, you should be right at home. Navigation controls are the standard you're used to.

#### Placing Pins
A pin can be placed on the map by single clicking anywhere on the map.

A note can be added to the pin by clicking on the pin to reveal the info popup. Clicking on the popup will
allow you to change the details of the pin. These notes will be included when your map is shared!

#### Removing Pins
To remove a pin, simply double click/tap on a pin. The pin will be removed.

#### Exporting Map
To share the map with others, click on the `export` icon in the upper right. Click on the `Copy to Clipboard`
button in the pop up, and the link will be saved to your active clipboard.

#### Reset Pins
To reset pins, click on the `reset` button in the upper right. All of your pins will be cleared.

### Limitations
There is a limit to the number of pins that you can share. The number of pins you can include depends on
whether or not you're including notes, and how long those notes are. Realistically, you should be able to
share 10 or so points, but this is **not** intended to be a way to plan out your life's mission.

## How It Works
Typically when sharing maps or some other thing with state, you need some way of maintaining that state. You
would normally stick that in a database and retrieve every time a request comes in. That's pretty trivial, and
common almost everywhere you go.

I wanted to approach this from a different angle. I wanted to be able to share maps without any kind of
database to hold state. In order to do this, I wanted to store the map state within the URL. While this is
obviously not the best way to approach this problem, it does work - at the cost of super ugly URLs.

The flow looks something like this:
1. User hits export on the map
2. A request is made to a serverless function with the map state in the payload
3. Serverless function will
    1. lzip the content of the payload
    2. base64 encode the payload
    3. return the encoded payload as a URL to UI
4. UI provides link to save the URL

And the flow for loading an exported map works the same, only in reverse:
1. User loads exported map
2. Query sent to serverless function
3. Serverless function will:
    1. base64 decode the payload
    2. lzip the content back to json
    3. return json payload
4. Browser loads the json payload and renders map + points


### Tech Used
* Svelte
* Bulma
* Bulma-toast
* Fontawesome

> Icons and favicon provided by [fontawesome](https://fontawesome.com/). No modifications were
made to their icons. [Licensing terms](https://fontawesome.com/license).

> All included maps belong to their respective creators. If you own the rights to a map and would
like it removed from the project, please send me a message. No map files are included in this
project, only links to where they are being publicly hosted.