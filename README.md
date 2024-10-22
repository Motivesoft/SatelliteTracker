# SatelliteTracker
Searches a feed of satellite data to display potential sightings

## Functionality
This is a command line app that issues calls against the [N2YO](https://www.n2yo.com/) web service, which provides satellite information.

When given a NORAD ID (satellite identification number), this app can return:
* The 'two line element' data for the satellite
* Information about potential sightings of the satellite from your location

My likely endgame is to develop this to notify me on a daily basis what satellite sightings might be possible from my location. I have not decided whether this will simply consist of textual output or maybe some sort of email notification, but feel free to fork the code and go your own way, and perhaps let me know what you come up with.

I do not intend to do anything pretty like mapping satellite information onto a map of the world.

## Code changes
Currently the code can perform a small collection of the N2YO API operations and is being build incrementally. If you want to experiment, check out the code and see how its built so that you can tweak it to your own aims.

### Development and Debugging
There is a boolean flag called `DEBUG` in the code that, when set to `true` will try and read values from a `json` text file rather than continually making calls to the web service, which has call limits. 

The idea is that during development, you would store the result of a successful call to the web API and then continue to develop the processing of that data by testing against the stored result before making any further API calls.

## Configuration
The following dotfiles are required to be in the same location as the application and their values configured as described below.

The precise meaning of the configuration values below can be worked out from the [N2YO API documentation](https://www.n2yo.com/api/).

### .env
Create a `.env` file with an API key value obtained from registering with [N2YO](https://www.n2yo.com/login).
```yaml
apiKey:
```

### .location
Create a `.location` file and complete with your location details as the 'observer' of any satellite sightings
```yaml
latitude:  # your latitude in decimal degrees format
longitude: # your longitude in decimal degrees format
altitude:  # your height in metres above sea level
```

### .preferences
Create a `.preferences` file and enter these general configuration settings
```yaml
days: 10               # number of days ahead to report, maximum 10
minimum_visibility: 60 # minimum seconds of visibility to consider, e.g. 60
```