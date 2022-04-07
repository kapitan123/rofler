Current.
1) Receive a message
2) Parse and get a link 
3) Make GetRequest, catch redirection to the page. 
4) Get new redirected result. Use it to download a video
3) Download from link
4) Convert to the right format


1) Links parser
	1) Parse all messages by regxp to check if the video is ticktok 
	/rofler realtime_parsing on
	/rofler realtime parsing off

	2) Parse the links from previous 10 messages /rofler prev 10
	3) Parse current /rofler https://tiktok.com/url/124124214

2) Owner 
	1) Adds Posted by rofler: @maypoldruha
	
3) Downloader
	1) Should download tik-tok file throught the save function
	2) Should limit downloads per day and per user if necessary
	3) Should use a separate account for dowloads
	4) I can youse an external service to download tiktoks
	
	https://github.com/siongui/tiktokgo
4) Convertor
	1) Should convert file to a suitable telegram format
		after download I need to convert it from mp4 to mpeg4 golang
		Video dimensions must be set to 480x320 (320x480 for vertical videos).
		H.264 and MPEG-4 should be used as the codec and container.
		https://core.telegram.org/blackberry/chat-media-send

5) Rofl counter
	1) Should count likes and emoji reactions on the post.
	2) Should count replyies to the post
		
6) Should store statistics
	1) PosterId, DatePosted, tiktok url, likescount
	2) /rofler top_rofl week shoud return the tiktokurl with top likes over the week
	3) Separate keker collection id\kekedOnDateTime\KekedOnTiktok\
	
7) Repostfinder OPTIONAL
	1) Should check if db contains tiktok uri

8) search videos by hashtag

// can be imported from the config
regexp to match ^https:\/\/vm\.tiktok\.com\/.*


Top rofler of the week is @maypoldruha (tiktoks posted:10, likes collected: 30, likes to tiktok ration 1/3)
Top keker @andreyivanov (keked: 10 times)

Add viewcount from the video.
