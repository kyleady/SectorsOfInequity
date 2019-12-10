from django.shortcuts import get_object_or_404, render, redirect
from django.forms import modelform_factory
from django.http import HttpResponse, Http404
import json
import os

class DefaultViews:
    def __init__(self, full_name, app, subapp, name, Model, custom=None):
        self.full_name = full_name
        self.title = name
        self.app = app
        self.subapp = subapp
        self.Model = Model
        self.custom = custom
        self.new_url = '{app}-{subapp}-{title}-new'.format(
            app=self.app,
            subapp=self.subapp,
            title=self.title
        )
        self.index_url = '{app}-{subapp}-{title}-index'.format(
            app=self.app,
            subapp=self.subapp,
            title=self.title
        )
        self.detail_url = '{app}-{subapp}-{title}-detail'.format(
            app=self.app,
            subapp=self.subapp,
            title=self.title
        )
        self.delete_url = '{app}-{subapp}-{title}-delete'.format(
            app=self.app,
            subapp=self.subapp,
            title=self.title
        )
        self.Form = modelform_factory(self.Model, fields='__all__')

        self.screaming_vortex_host = os.environ.get('SCREAMING_VORTEX_HOST')

    def index(self, request):
        template = 'index.html'
        context = {
            'new_url': self.new_url,
            'detail_url': self.detail_url,
            'full_name': self.full_name,
            'model_list': self.Model.objects.all()
        }

        return render(request, template, context)

    def detail(self, request, model_id):
        template = 'detail.html'
        model = get_object_or_404(self.Model, pk=model_id)
        if request.POST:
            form = self.Form(request.POST, instance=model)
            model = form.save()
            return redirect(self.index_url)

        context = {
            'model': model,
            'new_url': self.new_url,
            'detail_url': self.detail_url,
            'delete_url': self.delete_url,
            'full_name': self.full_name,
            'form': self.Form(instance=model)
        }

        return render(request, template, context)

    def new(self, request):
        template = 'detail.html'
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

    def delete(self, request, model_id):
        model = get_object_or_404(self.Model, pk=model_id)
        model.delete()
        return redirect(self.index_url)
