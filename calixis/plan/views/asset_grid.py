from django.shortcuts import get_object_or_404, render, redirect
from django import forms
import requests
import json

from .default import DefaultViews

class AssetGridViews(DefaultViews):
    def detail(self, request, model_id):
        template = 'asset_grid_detail.html'
        grid_model = get_object_or_404(self.Model, pk=model_id)

        node_models = self.custom['SubModels'][0].objects.filter(grid_id=grid_model.id)
        connection_models = self.custom['SubModels'][1].objects.filter(grid_id=grid_model.id)

        node_coords = []
        region_colors = {}
        max_x = 0
        max_y = 0
        for node_model in node_models:
            node_coords.append({
                "x": node_model.x,
                "y": node_model.y,
                "id": node_model.asset_id,
                # "region_id": node_model.region_id,
                "region_id": "PLACE_HOLDER",
            })

            if node_model.x > max_x:
                max_x = node_model.x
            if node_model.y > max_y:
                max_y = node_model.y


            #region_colors[node_model.region_id] = ''
            region_colors["PLACE_HOLDER"] = ''

        connection_coords = []
        for connection_model in connection_models:
            start_node = None
            end_node = None
            for node_model in node_models:
                if connection_model.start_id == node_model.id:
                    start_node = node_model
                if connection_model.end_id == node_model.id:
                    end_node = node_model
                if start_node is not None and end_node is not None:
                    break

            connection_coords.append({
                "start_x": start_node.x,
                "start_y": start_node.y,
                "end_x": end_node.x,
                "end_y": end_node.y,
                "id": connection_model.asset_id,
            })

        context = {
            'model': grid_model,
            'new_url': self.new_url,
            'detail_url': self.detail_url,
            'delete_url': self.delete_url,
            'full_name': self.full_name,
            'node_coords': node_coords,
            'connection_coords': connection_coords,
            'region_colors': region_colors,
            'canvas': {
                'max_x': None,
                'max_y': None,
                'offset_x': 10,
                'offset_y': 10,
                'spacing_x': 50,
                'spacing_y': 50,
                'radius': 10
            },
            'form': self.Form(instance=grid_model)
        }

        context['canvas']['max_x'] = 2 * context['canvas']['offset_x'] + max_x * context['canvas']['spacing_x']
        context['canvas']['max_y'] = 2 * context['canvas']['offset_y'] + max_y * context['canvas']['spacing_y']

        return render(request, template, context)
