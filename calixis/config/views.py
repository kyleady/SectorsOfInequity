from django.shortcuts import get_object_or_404, render, redirect
from django.forms import modelform_factory
from django.http import HttpResponse
import requests
import json
import os

class Views:
    def __init__(self, full_name, app, name, Model, SubModels=None):
        self.full_name = full_name
        self.title = name
        self.app = app
        self.Model = Model
        self.SubModels = SubModels
        self.new_url = '{app}-{title}-new'.format(
            app=self.app,
            title=self.title
        )
        self.index_url = '{app}-{title}-index'.format(
            app=self.app,
            title=self.title
        )
        self.detail_url = '{app}-{title}-detail'.format(
            app=self.app,
            title=self.title
        )
        self.test_url = '{app}-{title}-test'.format(
            app=self.app,
            title=self.title
        )
        self.Form = modelform_factory(self.Model, fields='__all__')

        self.screaming_vortex_host = os.environ.get('SCREAMING_VORTEX_HOST')

    def index(self, request):
        template = 'index.html'.format(
            title=self.title
        )
        context = {
            'new_url': self.new_url,
            'detail_url': self.detail_url,
            'full_name': self.full_name,
            'model_list': self.Model.objects.all()
        }

        return render(request, template, context)

    def detail(self, request, model_id):
        template = 'detail.html'.format(
            title=self.title
        )
        model = get_object_or_404(self.Model, pk=model_id)
        if request.POST:
            form = self.Form(request.POST, instance=model)
            model = form.save()
            return redirect(self.index_url)

        context = {
            'model': model,
            'new_url': self.new_url,
            'detail_url': self.detail_url,
            'full_name': self.full_name,
            'form': self.Form(instance=model)
        }

        return render(request, template, context)

    def sector_detail(self, request, model_id):
        template = 'sector_detail.html'
        sector_model = get_object_or_404(self.Model, pk=model_id)

        if request.POST:
            form = self.Form(request.POST, instance=sector_model)
            sector_model = form.save()
            return redirect(self.index_url)

        system_models = self.SubModels[0].objects.filter(sector_id=sector_model.id)
        route_models = self.SubModels[1].objects.filter(start__in=system_models)

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
        template = 'detail.html'.format(
            title=self.title
        )
        if request.POST:
            form = self.Form(request.POST)
            model = form.save()
            return redirect(self.index_url)

        context = {
            'full_name': self.full_name,
            'new_url': self.new_url,
            'detail_url': self.detail_url,
            'form': self.Form()
        }

        return render(request, template, context)

    def test(self, request, model_id):
        template = 'test.html'.format(
            title=self.title
        )

        screaming_vortex_url = 'http://{host}/{path}'.format(
            host=self.screaming_vortex_host,
            path=self.title
        )

        template_id = None
        model = get_object_or_404(self.Model, pk=model_id)

        screaming_vortex_json = {
            'config_id': model_id,
            'skeleton_id': None,
        }
        screaming_vortex_response = requests.post(
            screaming_vortex_url,
            data=json.dumps(screaming_vortex_json),
        )
        screaming_vortex_response.raise_for_status()
        return HttpResponse(screaming_vortex_response.text)
