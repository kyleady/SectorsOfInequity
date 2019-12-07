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
  var canvas = document.getElementById('canvas');
  if (canvas.getContext) {
    var ctx = canvas.getContext('2d');
    var system_coords = {{ system_coords|safe }}
    var route_coords = {{ route_coords|safe }}
    var region_colors = {{ region_colors|safe }}
    var radius = {{ canvas.radius }}

    var spectrum = createSpectrum(Object.keys(region_colors).length);
    console.log(spectrum)
    var color_i = 0
    for (var region in region_colors) {
      region_colors[region] = spectrum[color_i]
      color_i++
    }

    for (var system of system_coords) {
      ctx.fillStyle = region_colors[system.region_id]
      console.log(ctx.fillStyle)
      system_at = gridToPixels(system.x, system.y)
      ctx.beginPath();
      ctx.arc(system_at.x, system_at.y, radius, 0, 2 * Math.PI);
      ctx.fill()
    }

    for (var route of route_coords) {
      start_at = gridToPixels(route.start_x, route.start_y)
      end_at = gridToPixels(route.end_x, route.end_y)
      ctx.beginPath();
      ctx.moveTo(start_at.x, start_at.y)
      ctx.lineTo(end_at.x, end_at.y)
      ctx.stroke()
    }
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
