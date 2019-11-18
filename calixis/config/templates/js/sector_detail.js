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
    var radius = {{ canvas.radius }}


    for (var system of system_coords) {
      system_at = gridToPixels(system.x, system.y)
      ctx.beginPath();
      ctx.arc(system_at.x, system_at.y, radius, 0, 2 * Math.PI);
      ctx.stroke()
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
