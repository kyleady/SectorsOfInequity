from django.shortcuts import get_object_or_404, render, redirect
from django import forms
import requests
import json

from .default import DefaultViews



class GridSectorViews(DefaultViews):
    def detail(self, request, model_id):
        template = 'sector_detail.html'
        sector_model = get_object_or_404(self.Model, pk=model_id)

        if request.POST:
            form = self.Form(request.POST, instance=sector_model)
            sector_model = form.save()
            return redirect(self.index_url)

        system_models = self.custom['SubModels'][0].objects.filter(sector_id=sector_model.id)
        route_models = self.custom['SubModels'][1].objects.filter(start__in=system_models)

        system_coords = []
        region_colors = {}
        max_x = 0
        max_y = 0
        for system_model in system_models:
            system_coords.append({
                "x": system_model.x,
                "y": system_model.y,
                "region_id": system_model.region_id,
            })

            if system_model.x > max_x:
                max_x = system_model.x
            if system_model.y > max_y:
                max_y = system_model.y

            region_colors[system_model.region_id] = ''

        route_coords = []
        for route_model in route_models:
            start_system = None
            end_system = None
            for system_model in system_models:
                if route_model.start_id == system_model.id:
                    start_system = system_model
                if route_model.end_id == system_model.id:
                    end_system = system_model
                if start_system is not None and end_system is not None:
                    break

            route_coords.append({
                "start_x": start_system.x,
                "start_y": start_system.y,
                "end_x": end_system.x,
                "end_y": end_system.y,
            })

        context = {
            'model': sector_model,
            'new_url': self.new_url,
            'detail_url': self.detail_url,
            'delete_url': self.delete_url,
            'full_name': self.full_name,
            'system_coords': system_coords,
            'route_coords': route_coords,
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
            'form': self.Form(instance=sector_model)
        }

        context['canvas']['max_x'] = 2 * context['canvas']['offset_x'] + max_x * context['canvas']['spacing_x']
        context['canvas']['max_y'] = 2 * context['canvas']['offset_y'] + max_y * context['canvas']['spacing_y']

        return render(request, template, context)

    def new(self, request):
        allGridConfig = self.custom['Grid'].objects
        class GenerateSectorConfigForm(forms.Form):
            config = forms.ModelChoiceField(queryset=allGridConfig)
        if request.POST:
            form = GenerateSectorConfigForm(request.POST)
            form.is_valid()
            config_id = form.cleaned_data['config'].id
            screaming_vortex_url = 'http://{host}/{path}'.format(
                host=self.screaming_vortex_host,
                path='grid'
            )
            screaming_vortex_json = {
                'config_id': config_id,
                'skeleton_id': None,
            }
            screaming_vortex_response = requests.post(
                screaming_vortex_url,
                data=json.dumps(screaming_vortex_json),
            )
            screaming_vortex_response.raise_for_status()
            return redirect(self.index_url)
        else:
            template = 'detail.html'
            context = {
                'full_name': self.full_name,
                'new_url': self.new_url,
                'detail_url': self.detail_url,
                'form': GenerateSectorConfigForm()
            }

            return render(request, template, context)
