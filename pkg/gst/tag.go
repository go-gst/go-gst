package gst

// TAG_ALBUM (GST_TAG_ALBUM): album containing this data (string)
//
// The album name as it should be displayed, e.g. 'The Jazz Guitar'.
const TAG_ALBUM = "album"

// TAG_ALBUM_ARTIST (GST_TAG_ALBUM_ARTIST): artist of the entire album, as it
// should be displayed.
const TAG_ALBUM_ARTIST = "album-artist"

// TAG_ALBUM_ARTIST_SORTNAME (GST_TAG_ALBUM_ARTIST_SORTNAME): artist of the
// entire album, as it should be sorted.
const TAG_ALBUM_ARTIST_SORTNAME = "album-artist-sortname"

// TAG_ALBUM_GAIN (GST_TAG_ALBUM_GAIN): album gain in db (double).
const TAG_ALBUM_GAIN = "replaygain-album-gain"

// TAG_ALBUM_PEAK (GST_TAG_ALBUM_PEAK): peak of the album (double).
const TAG_ALBUM_PEAK = "replaygain-album-peak"

// TAG_ALBUM_SORTNAME (GST_TAG_ALBUM_SORTNAME): album containing this data,
// as used for sorting (string)
//
// The album name as it should be sorted, e.g. 'Jazz Guitar, The'.
const TAG_ALBUM_SORTNAME = "album-sortname"

// TAG_ALBUM_VOLUME_COUNT (GST_TAG_ALBUM_VOLUME_COUNT): count of discs inside
// collection this disc belongs to (unsigned integer).
const TAG_ALBUM_VOLUME_COUNT = "album-disc-count"

// TAG_ALBUM_VOLUME_NUMBER (GST_TAG_ALBUM_VOLUME_NUMBER): disc number inside a
// collection (unsigned integer).
const TAG_ALBUM_VOLUME_NUMBER = "album-disc-number"

// TAG_APPLICATION_DATA (GST_TAG_APPLICATION_DATA): arbitrary application data
// (sample)
//
// Some formats allow applications to add their own arbitrary data into files.
// This data is application dependent.
const TAG_APPLICATION_DATA = "application-data"

// TAG_APPLICATION_NAME (GST_TAG_APPLICATION_NAME): name of the application used
// to create the media (string).
const TAG_APPLICATION_NAME = "application-name"

// TAG_ARTIST (GST_TAG_ARTIST): person(s) responsible for the recording (string)
//
// The artist name as it should be displayed, e.g. 'Jimi Hendrix' or 'The Guitar
// Heroes'.
const TAG_ARTIST = "artist"

// TAG_ARTIST_SORTNAME (GST_TAG_ARTIST_SORTNAME): person(s) responsible for the
// recording, as used for sorting (string)
//
// The artist name as it should be sorted, e.g. 'Hendrix, Jimi' or 'Guitar
// Heroes, The'.
const TAG_ARTIST_SORTNAME = "artist-sortname"

// TAG_ATTACHMENT (GST_TAG_ATTACHMENT): generic file attachment (sample) (sample
// taglist should specify the content type and if possible set "filename" to the
// file name of the attachment).
const TAG_ATTACHMENT = "attachment"

// TAG_AUDIO_CODEC (GST_TAG_AUDIO_CODEC): codec the audio data is stored in
// (string).
const TAG_AUDIO_CODEC = "audio-codec"

// TAG_BEATS_PER_MINUTE (GST_TAG_BEATS_PER_MINUTE): number of beats per minute
// in audio (double).
const TAG_BEATS_PER_MINUTE = "beats-per-minute"

// TAG_BITRATE (GST_TAG_BITRATE): exact or average bitrate in bits/s (unsigned
// integer).
const TAG_BITRATE = "bitrate"

// TAG_CODEC (GST_TAG_CODEC): codec the data is stored in (string).
const TAG_CODEC = "codec"

// TAG_COMMENT (GST_TAG_COMMENT): free text commenting the data (string).
const TAG_COMMENT = "comment"

// TAG_COMPOSER (GST_TAG_COMPOSER): person(s) who composed the recording
// (string).
const TAG_COMPOSER = "composer"

// TAG_COMPOSER_SORTNAME (GST_TAG_COMPOSER_SORTNAME) composer's name, used for
// sorting (string).
const TAG_COMPOSER_SORTNAME = "composer-sortname"

// TAG_CONDUCTOR (GST_TAG_CONDUCTOR): conductor/performer refinement (string).
const TAG_CONDUCTOR = "conductor"

// TAG_CONTACT (GST_TAG_CONTACT): contact information (string).
const TAG_CONTACT = "contact"

// TAG_CONTAINER_FORMAT (GST_TAG_CONTAINER_FORMAT): container format the data is
// stored in (string).
const TAG_CONTAINER_FORMAT = "container-format"

// TAG_CONTAINER_SPECIFIC_TRACK_ID (GST_TAG_CONTAINER_SPECIFIC_TRACK_ID):
// unique identifier for the audio, video or text track this tag is associated
// with. The mappings for several container formats are defined in the
// [Sourcing In-band Media Resource Tracks from Media Containers into HTML
// specification](https://dev.w3.org/html5/html-sourcing-inband-tracks/).
const TAG_CONTAINER_SPECIFIC_TRACK_ID = "container-specific-track-id"

// TAG_COPYRIGHT (GST_TAG_COPYRIGHT): copyright notice of the data (string).
const TAG_COPYRIGHT = "copyright"

// TAG_COPYRIGHT_URI (GST_TAG_COPYRIGHT_URI): URI to location where copyright
// details can be found (string).
const TAG_COPYRIGHT_URI = "copyright-uri"

// TAG_DATE (GST_TAG_DATE): date the data was created (#GDate structure).
const TAG_DATE = "date"

// TAG_DATE_TIME (GST_TAG_DATE_TIME): date and time the data was created
// (DateTime structure).
const TAG_DATE_TIME = "datetime"

// TAG_DESCRIPTION (GST_TAG_DESCRIPTION): short text describing the content of
// the data (string).
const TAG_DESCRIPTION = "description"

// TAG_DEVICE_MANUFACTURER (GST_TAG_DEVICE_MANUFACTURER): manufacturer of the
// device used to create the media (string).
const TAG_DEVICE_MANUFACTURER = "device-manufacturer"

// TAG_DEVICE_MODEL (GST_TAG_DEVICE_MODEL): model of the device used to create
// the media (string).
const TAG_DEVICE_MODEL = "device-model"

// TAG_DURATION (GST_TAG_DURATION): length in GStreamer time units (nanoseconds)
// (unsigned 64-bit integer).
const TAG_DURATION = "duration"

// TAG_ENCODED_BY (GST_TAG_ENCODED_BY): name of the person or organisation
// that encoded the file. May contain a copyright message if the person or
// organisation also holds the copyright (string)
//
// Note: do not use this field to describe the encoding application. Use
// T_TAG_APPLICATION_NAME or T_TAG_COMMENT for that.
const TAG_ENCODED_BY = "encoded-by"

// TAG_ENCODER (GST_TAG_ENCODER): encoder used to encode this stream (string).
const TAG_ENCODER = "encoder"

// TAG_ENCODER_VERSION (GST_TAG_ENCODER_VERSION): version of the encoder used to
// encode this stream (unsigned integer).
const TAG_ENCODER_VERSION = "encoder-version"

// TAG_EXTENDED_COMMENT (GST_TAG_EXTENDED_COMMENT): key/value text commenting
// the data (string)
//
// Must be in the form of 'key=comment' or 'key[lc]=comment' where 'lc' is an
// ISO-639 language code.
//
// This tag is used for unknown Vorbis comment tags, unknown APE tags and
// certain ID3v2 comment fields.
const TAG_EXTENDED_COMMENT = "extended-comment"

// TAG_GENRE (GST_TAG_GENRE): genre this data belongs to (string).
const TAG_GENRE = "genre"

// TAG_GEO_LOCATION_CAPTURE_DIRECTION (GST_TAG_GEO_LOCATION_CAPTURE_DIRECTION)
// indicates the direction the device is pointing to when capturing a media.
// It is represented as degrees in floating point representation, 0 means the
// geographic north, and increases clockwise (double from 0 to 360)
//
// See also T_TAG_GEO_LOCATION_MOVEMENT_DIRECTION.
const TAG_GEO_LOCATION_CAPTURE_DIRECTION = "geo-location-capture-direction"

// TAG_GEO_LOCATION_CITY (GST_TAG_GEO_LOCATION_CITY): city (english name) where
// the media has been produced (string).
const TAG_GEO_LOCATION_CITY = "geo-location-city"

// TAG_GEO_LOCATION_COUNTRY (GST_TAG_GEO_LOCATION_COUNTRY): country (english
// name) where the media has been produced (string).
const TAG_GEO_LOCATION_COUNTRY = "geo-location-country"

// TAG_GEO_LOCATION_ELEVATION (GST_TAG_GEO_LOCATION_ELEVATION): geo elevation
// of where the media has been recorded or produced in meters according to WGS84
// (zero is average sea level) (double).
const TAG_GEO_LOCATION_ELEVATION = "geo-location-elevation"

// TAG_GEO_LOCATION_HORIZONTAL_ERROR (GST_TAG_GEO_LOCATION_HORIZONTAL_ERROR)
// represents the expected error on the horizontal positioning in meters
// (double).
const TAG_GEO_LOCATION_HORIZONTAL_ERROR = "geo-location-horizontal-error"

// TAG_GEO_LOCATION_LATITUDE (GST_TAG_GEO_LOCATION_LATITUDE): geo latitude
// location of where the media has been recorded or produced in degrees
// according to WGS84 (zero at the equator, negative values for southern
// latitudes) (double).
const TAG_GEO_LOCATION_LATITUDE = "geo-location-latitude"

// TAG_GEO_LOCATION_LONGITUDE (GST_TAG_GEO_LOCATION_LONGITUDE): geo longitude
// location of where the media has been recorded or produced in degrees
// according to WGS84 (zero at the prime meridian in Greenwich/UK, negative
// values for western longitudes). (double).
const TAG_GEO_LOCATION_LONGITUDE = "geo-location-longitude"

// TAG_GEO_LOCATION_MOVEMENT_DIRECTION (GST_TAG_GEO_LOCATION_MOVEMENT_DIRECTION)
// indicates the movement direction of the device performing the capture of
// a media. It is represented as degrees in floating point representation,
// 0 means the geographic north, and increases clockwise (double from 0 to 360)
//
// See also T_TAG_GEO_LOCATION_CAPTURE_DIRECTION.
const TAG_GEO_LOCATION_MOVEMENT_DIRECTION = "geo-location-movement-direction"

// TAG_GEO_LOCATION_MOVEMENT_SPEED (GST_TAG_GEO_LOCATION_MOVEMENT_SPEED): speed
// of the capturing device when performing the capture. Represented in m/s.
// (double)
//
// See also T_TAG_GEO_LOCATION_MOVEMENT_DIRECTION.
const TAG_GEO_LOCATION_MOVEMENT_SPEED = "geo-location-movement-speed"

// TAG_GEO_LOCATION_NAME (GST_TAG_GEO_LOCATION_NAME): human readable descriptive
// location of where the media has been recorded or produced. (string).
const TAG_GEO_LOCATION_NAME = "geo-location-name"

// TAG_GEO_LOCATION_SUBLOCATION (GST_TAG_GEO_LOCATION_SUBLOCATION): location
// 'smaller' than GST_TAG_GEO_LOCATION_CITY that specifies better where the
// media has been produced. (e.g. the neighborhood) (string).
//
// This tag has been added as this is how it is handled/named in XMP's
// Iptc4xmpcore schema.
const TAG_GEO_LOCATION_SUBLOCATION = "geo-location-sublocation"

// TAG_GROUPING (GST_TAG_GROUPING) groups together media that are related
// and spans multiple tracks. An example are multiple pieces of a concerto.
// (string).
const TAG_GROUPING = "grouping"

// TAG_HOMEPAGE (GST_TAG_HOMEPAGE): homepage for this media (i.e. artist or
// movie homepage) (string).
const TAG_HOMEPAGE = "homepage"

// TAG_IMAGE (GST_TAG_IMAGE): image (sample) (sample taglist should specify the
// content type and preferably also set "image-type" field as GstTagImageType).
const TAG_IMAGE = "image"

// TAG_IMAGE_ORIENTATION (GST_TAG_IMAGE_ORIENTATION) represents the
// 'Orientation' tag from EXIF. Defines how the image should be rotated and
// mirrored for display. (string)
//
// This tag has a predefined set of allowed values: "rotate-0" "rotate-90"
// "rotate-180" "rotate-270" "flip-rotate-0" "flip-rotate-90" "flip-rotate-180"
// "flip-rotate-270"
//
// The naming is adopted according to a possible transformation to perform on
// the image to fix its orientation, obviously equivalent operations will yield
// the same result.
//
// Rotations indicated by the values are in clockwise direction and 'flip' means
// an horizontal mirroring.
const TAG_IMAGE_ORIENTATION = "image-orientation"

// TAG_INTERPRETED_BY (GST_TAG_INTERPRETED_BY): information about the people
// behind a remix and similar interpretations of another existing piece
// (string).
const TAG_INTERPRETED_BY = "interpreted-by"

// TAG_ISRC (GST_TAG_ISRC): international Standard Recording Code - see
// http://www.ifpi.org/isrc/ (string).
const TAG_ISRC = "isrc"

// TAG_KEYWORDS (GST_TAG_KEYWORDS): comma separated keywords describing the
// content (string).
const TAG_KEYWORDS = "keywords"

// TAG_LANGUAGE_CODE (GST_TAG_LANGUAGE_CODE): ISO-639-2 or ISO-639-1 code for
// the language the content is in (string)
//
// There is utility API in libgsttag in gst-plugins-base to obtain a translated
// language name from the language code: gst_tag_get_language_name().
const TAG_LANGUAGE_CODE = "language-code"

// TAG_LANGUAGE_NAME (GST_TAG_LANGUAGE_NAME): name of the language the content
// is in (string)
//
// Free-form name of the language the content is in, if a language code is
// not available. This tag should not be set in addition to a language code.
// It is undefined what language or locale the language name is in.
const TAG_LANGUAGE_NAME = "language-name"

// TAG_LICENSE (GST_TAG_LICENSE): license of data (string).
const TAG_LICENSE = "license"

// TAG_LICENSE_URI (GST_TAG_LICENSE_URI): URI to location where license details
// can be found (string).
const TAG_LICENSE_URI = "license-uri"

// TAG_LOCATION (GST_TAG_LOCATION): origin of media as a URI (location, where
// the original of the file or stream is hosted) (string).
const TAG_LOCATION = "location"

// TAG_LYRICS (GST_TAG_LYRICS) lyrics of the media (string).
const TAG_LYRICS = "lyrics"

// TAG_MAXIMUM_BITRATE (GST_TAG_MAXIMUM_BITRATE): maximum bitrate in bits/s
// (unsigned integer).
const TAG_MAXIMUM_BITRATE = "maximum-bitrate"

// TAG_MIDI_BASE_NOTE (GST_TAG_MIDI_BASE_NOTE): Midi note number
// (http://en.wikipedia.org/wiki/Note#Note_designation_in_accordance_with_octave_name)
// of the audio track. This is useful for sample instruments and in particular
// for multi-samples.
const TAG_MIDI_BASE_NOTE = "midi-base-note"

// TAG_MINIMUM_BITRATE (GST_TAG_MINIMUM_BITRATE): minimum bitrate in bits/s
// (unsigned integer).
const TAG_MINIMUM_BITRATE = "minimum-bitrate"

// TAG_NOMINAL_BITRATE (GST_TAG_NOMINAL_BITRATE): nominal bitrate in bits/s
// (unsigned integer). The actual bitrate might be different from this target
// bitrate.
const TAG_NOMINAL_BITRATE = "nominal-bitrate"

// TAG_ORGANIZATION (GST_TAG_ORGANIZATION): organization (string).
const TAG_ORGANIZATION = "organization"

// TAG_PERFORMER (GST_TAG_PERFORMER): person(s) performing (string).
const TAG_PERFORMER = "performer"

// TAG_PREVIEW_IMAGE (GST_TAG_PREVIEW_IMAGE): image that is meant for preview
// purposes, e.g. small icon-sized version (sample) (sample taglist should
// specify the content type).
const TAG_PREVIEW_IMAGE = "preview-image"

// TAG_PRIVATE_DATA (GST_TAG_PRIVATE_DATA): any private data that may be
// contained in tags (sample).
//
// It is represented by Sample in which Buffer contains the binary data and the
// sample's info Structure may contain any extra information that identifies the
// origin or meaning of the data.
//
// Private frames in ID3v2 tags ('PRIV' frames) will be represented using
// this tag, in which case the GstStructure will be named "ID3PrivateFrame"
// and contain a field named "owner" of type string which contains the
// owner-identification string from the tag.
const TAG_PRIVATE_DATA = "private-data"

// TAG_PUBLISHER (GST_TAG_PUBLISHER): name of the label or publisher (string).
const TAG_PUBLISHER = "publisher"

// TAG_REFERENCE_LEVEL (GST_TAG_REFERENCE_LEVEL): reference level of track and
// album gain values (double).
const TAG_REFERENCE_LEVEL = "replaygain-reference-level"

// TAG_SERIAL (GST_TAG_SERIAL): serial number of track (unsigned integer).
const TAG_SERIAL = "serial"

// TAG_SHOW_EPISODE_NUMBER (GST_TAG_SHOW_EPISODE_NUMBER): number of the episode
// within a season/show (unsigned integer).
const TAG_SHOW_EPISODE_NUMBER = "show-episode-number"

// TAG_SHOW_NAME (GST_TAG_SHOW_NAME): name of the show, used for displaying
// (string).
const TAG_SHOW_NAME = "show-name"

// TAG_SHOW_SEASON_NUMBER (GST_TAG_SHOW_SEASON_NUMBER): number of the season of
// a show/series (unsigned integer).
const TAG_SHOW_SEASON_NUMBER = "show-season-number"

// TAG_SHOW_SORTNAME (GST_TAG_SHOW_SORTNAME): name of the show, used for sorting
// (string).
const TAG_SHOW_SORTNAME = "show-sortname"

// TAG_SUBTITLE_CODEC (GST_TAG_SUBTITLE_CODEC): codec/format the subtitle data
// is stored in (string).
const TAG_SUBTITLE_CODEC = "subtitle-codec"

// TAG_TITLE (GST_TAG_TITLE): commonly used title (string)
//
// The title as it should be displayed, e.g. 'The Doll House'.
const TAG_TITLE = "title"

// TAG_TITLE_SORTNAME (GST_TAG_TITLE_SORTNAME): commonly used title, as used for
// sorting (string)
//
// The title as it should be sorted, e.g. 'Doll House, The'.
const TAG_TITLE_SORTNAME = "title-sortname"

// TAG_TRACK_COUNT (GST_TAG_TRACK_COUNT): count of tracks inside collection this
// track belongs to (unsigned integer).
const TAG_TRACK_COUNT = "track-count"

// TAG_TRACK_GAIN (GST_TAG_TRACK_GAIN): track gain in db (double).
const TAG_TRACK_GAIN = "replaygain-track-gain"

// TAG_TRACK_NUMBER (GST_TAG_TRACK_NUMBER): track number inside a collection
// (unsigned integer).
const TAG_TRACK_NUMBER = "track-number"

// TAG_TRACK_PEAK (GST_TAG_TRACK_PEAK): peak of the track (double).
const TAG_TRACK_PEAK = "replaygain-track-peak"

// TAG_USER_RATING (GST_TAG_USER_RATING): rating attributed by a person (likely
// the application user). The higher the value, the more the user likes this
// media (unsigned int from 0 to 100).
const TAG_USER_RATING = "user-rating"

// TAG_VERSION (GST_TAG_VERSION): version of this data (string).
const TAG_VERSION = "version"

// TAG_VIDEO_CODEC (GST_TAG_VIDEO_CODEC): codec the video data is stored in
// (string).
const TAG_VIDEO_CODEC = "video-codec"
