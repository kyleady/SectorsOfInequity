function gridToPixels(grid_x, grid_y) {
  var offset_x = {{ canvas.offset_x }}
  var offset_y = {{ canvas.offset_y }}
  var spacing_x = {{ canvas.spacing_x }}
  var spacing_y = {{ canvas.spacing_y }}

  return {
    "x": offset_x + (grid_x * spacing_x),
    "y": offset_y + (grid_y * spacing_y),
  }
}

function draw() {
  var Links = new Array();
  var hoverLink = "";
  var canvas = document.getElementById('canvas');
  if (canvas.getContext) {
    var ctx = canvas.getContext('2d');
    var node_coords = {{ node_coords|safe }}
    var connection_coords = {{ connection_coords|safe }}
    var region_colors = {{ region_colors|safe }}
    var radius = {{ canvas.radius }}

    var spectrum = createSpectrum(Object.keys(region_colors).length);
    console.log(spectrum)
    var color_i = 0
    for (var region in region_colors) {
      region_colors[region] = spectrum[color_i]
      color_i++
    }

    for (var node of node_coords) {
      ctx.fillStyle = region_colors[node.region_id]
      console.log(ctx.fillStyle)
      node_at = gridToPixels(node.x, node.y)
      ctx.beginPath();
      ctx.arc(node_at.x, node_at.y, radius, 0, 2 * Math.PI);
      ctx.fill()

      Links.push({
        X: node_at.x - radius,
        Y: node_at.y - radius,
        Width: 2 * radius,
        Height: 2 * radius,
        Href: "/plan/asset/asset/" + node.id
      })
    }

    ctx.fillStyle = "#000000"
    for (var connection of connection_coords) {
      start_at = gridToPixels(connection.start_x, connection.start_y)
      end_at = gridToPixels(connection.end_x, connection.end_y)
      midpoint_at = {
        "x": (start_at.x + end_at.x)/2,
        "y": (start_at.y + end_at.y)/2,
      }

      ctx.beginPath();
      ctx.moveTo(start_at.x, start_at.y)
      ctx.lineTo(end_at.x, end_at.y)
      ctx.stroke()

      ctx.beginPath()
      ctx.arc(midpoint_at.x, midpoint_at.y, radius / 4, 0, 2 * Math.PI)
      ctx.fill()

      Links.push({
        X: midpoint_at.x - radius / 4,
        Y: midpoint_at.y - radius / 4,
        Width: radius / 2,
        Height: radius / 2,
        Href: "/plan/asset/asset/" + connection.id
      })
    }

    // Link hover
    function on_mousemove (ev) {
        var x, y;

        // Get the mouse position relative to the canvas element
        if (ev.layerX || ev.layerX == 0) { // For Firefox
            x = ev.layerX;
            y = ev.layerY;
        }

        // Link hover
        for (var i = 0; i < Links.length; i++) {
            var link = Links[i];

            // Check if cursor is in the link area
            if (x >= link.X && x <= (link.X + link.Width)
              && y >= link.Y && y <= (link.Y + link.Height)){
                document.body.style.cursor = "pointer";
                hoverLink = link.Href;
                break;
            } else {
                document.body.style.cursor = "";
                hoverLink = "";
            }
        };
    }

    // Link click
    function on_click(e) {
        if (hoverLink) {
            window.location = hoverLink;
        }
    }

    canvas.addEventListener("mousemove", on_mousemove, false);
    canvas.addEventListener("click", on_click, false);
  }
}

//START: https://gist.github.com/ibrechin/2489005
hslToRgb = function(_h, s, l) {
	var h = Math.min(_h, 359)/60;

	var c = (1-Math.abs((2*l)-1))*s;
	var x = c*(1-Math.abs((h % 2)-1));
	var m = l - (0.5*c);

	var r = m, g = m, b = m;

	if (h < 1) {
		r += c, g = +x, b += 0;
	} else if (h < 2) {
		r += x, g += c, b += 0;
	} else if (h < 3) {
		r += 0, g += c, b += x;
	} else if (h < 4) {
		r += 0, g += x, b += c;
	} else if (h < 5) {
		r += x, g += 0, b += c;
	} else if (h < 6) {
		r += c, g += 0, b += x;
	} else {
		r = 0, g = 0, b = 0;
	}

	return 'rgb(' + Math.floor(r*255) + ', ' + Math.floor(g*255) + ', ' + Math.floor(b*255) + ')';
}

createSpectrum = function(length) {
	var colors = [];
	// 270 because we don't want the spectrum to circle back
	var step = 270/length;
	for (var i = 1; i <= length; i++) {
		var color = hslToRgb((i)*step, 0.5, 0.5);
		colors.push(color);
	}

	return colors;
}
//END: https://gist.github.com/ibrechin/2489005
