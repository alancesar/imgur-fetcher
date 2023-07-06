package testdata

const ImgurImageResponse = `
{
  "data": {
    "id": "some-image-id",
	"description": "Some image description",
    "title": "Some image title",
    "type": "image\/jpeg",
    "width": 1080,
    "height": 1920,
    "link": "https:\/\/i.imgur.com\/some-image.jpg"
  },
  "success": true,
  "status": 200
}
`

const ImgurVideoResponse = `
{
  "data": {
    "id": "some-video-id",
    "title": "Some video title",
    "description": "Some video description",
    "type": "video\/mp4",
    "width": 1920,
    "height": 1080,
    "link": "https:\/\/i.imgur.com\/some-video.mp4",
    "mp4": "https:\/\/i.imgur.com\/some-video.mp4",
    "gifv": "https:\/\/i.imgur.com\/some-video.gifv",
    "hls": "https:\/\/i.imgur.com\/some-video.m3u8"
  },
  "success": true,
  "status": 200
}
`

const ImgurAlbumResponse = `
{
  "data": {
    "id": "some-album",
    "title": "Some album title",
    "description": "Some album description",
    "link": "https:\/\/imgur.com\/a\/some-album",
    "is_album": true,
    "images": [
      {
        "id": "some-image-id-1",
        "title": "Some image #1 title",
        "description": "Some image #1 description",
        "type": "image\/jpeg",
        "width": 1080,
        "height": 1920,
        "link": "https:\/\/i.imgur.com\/some-image-1.jpg"
      },
      {
        "id": "some-image-id-2",
        "title": "Some image #2 title",
        "description": "Some image #2 description",
        "type": "image\/jpeg",
        "width": 720,
        "height": 1280,
        "link": "https:\/\/i.imgur.com\/some-image-2.jpg"
      }
    ]
  },
  "success": true,
  "status": 200
}
`
