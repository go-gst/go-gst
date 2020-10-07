package video

/*
#include <gst/video/video.h>
*/
import "C"

// TileMode is an enum value describing the available tiling modes.
type TileMode int

// Type castings
const (
	TileModeUnknown   TileMode = C.GST_VIDEO_TILE_MODE_UNKNOWN    // (0) – Unknown or unset tile mode
	TileModeZFlipZ2X2 TileMode = C.GST_VIDEO_TILE_MODE_ZFLIPZ_2X2 // (65536) – Every four adjacent blocks - two horizontally and two vertically are grouped together and are located in memory in Z or flipped Z order. In case of odd rows, the last row of blocks is arranged in linear order.
	TileModeLinear    TileMode = C.GST_VIDEO_TILE_MODE_LINEAR     // (131072) – Tiles are in row order.
)

// TileType is an enum value describing the most common tiling types.
type TileType int

// Type castings
const (
	TileTypeIndexed TileType = C.GST_VIDEO_TILE_TYPE_INDEXED // (0) – Tiles are indexed. Use gst_video_tile_get_index () to retrieve the tile at the requested coordinates.
)
