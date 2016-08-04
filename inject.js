var pointDiv = "<div id=\"pztacker-pointer\" style=\"width: 8px; height: 8px; position: absolute; left: calc(50% - 4px); top: calc(50% - 4px); background-color: rgba(0,0,0,0.2); z-index: 100; border-radius: 5px; border: 1px solid rgba(255, 0, 0, 0.7);\"></div>"

function getUpdate() {
	var connection = new WebSocket('ws://127.0.0.1:%PORT%/websocket')
	
	// Add Pointer if not existend
	if($('#pztacker-pointer').length == 0) {
		$(".openseadragon-container").append(pointDiv);
	}

	// Websocket Packet handling
	connection.onmessage = function (e) {
	  obj = JSON.parse(e.data);
	  
	  if (obj.Name == "Position") {
			c = tileToPixel(parseFloat(obj.Data.X), parseFloat(obj.Data.Y));
			viewer.viewport.panTo(viewer.viewport.imageToViewportCoordinates(c.x, c.y));
	  } else if (obj.Name == "ZoomInPress") {
			viewer.viewport.zoomTo(viewer.viewport.getZoom() + 25, null, true);
	  } else if (obj.Name == "ZoomOutPress") {
			viewer.viewport.zoomTo(viewer.viewport.getZoom() - 25, null, true);
	  } else if (obj.Name == "UpPress") {
			go(0, -(0.0015 * (100/viewer.viewport.getZoom())));
	  } else if (obj.Name == "DownPress") {
			go(0, (0.0015 * (100/viewer.viewport.getZoom())));
	  } else if (obj.Name == "LeftPress") {
			go(-(0.0015 * (100/viewer.viewport.getZoom())), 0);
	  } else if (obj.Name == "RightPress") {
			go((0.0015 * (100/viewer.viewport.getZoom())), 0);
	  }
	};
}

// Moves the map
function go(dX, dY) {
	c = viewer.viewport.getCenter();
	c.x += dX;
	c.y += dY;
	viewer.viewport.panTo(c);
}

getUpdate();