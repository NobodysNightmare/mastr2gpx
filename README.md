# MaStR2GPX

This is a tool that can be used to convert a full dump ("Gesamtdatenauszug") of the german
[Marktstammdatenregister](https://www.marktstammdatenregister.de) into a GPX file:

```
$ mastr2gpx Gesamtdatenexport_20241012_24.2/
```

This is intended to be a mapping help, e.g. when correlating official data to things you see on [OpenStreetMap](https://www.openstreetmap.org)
or on a satellite image.

## Filtering data

The tool allows to filter data based on the postal code or a bounding box by providing command-line parameters, e.g.

```
$ mastr2gpx --postal-code 12345 Gesamtdatenexport_20241012_24.2/
```

```
$ mastr2gpx --bbox 12.6332,52.5317,12.8725,52.6316 Gesamtdatenexport_20241012_24.2/
```

## Accuracy warning

I built this to help my own mapping activities on OpenStreetMap. However, while using it, I realized that the official
geo coordinates are often only appoximated and may be off by hundreds of meters. I had instances of coordinates matching
a solar park on the map, but the data didn't align (wrong panel orientation, too few panels).

So make sure to perform some sanity checks before using the data on OSM.

## Disclaimer & Acknowledgement

Just in case this wasn't obvious already: This tool is **not official** in any kind or form. While it ingests data from the Marktstammdatenregister,
it is not provided by the Bundesnetzagentur or the Marktstammdatenregister itself. It just works using the openly available data offered by
the register.

Thanks for making this kind of data publicly available! ❤️
