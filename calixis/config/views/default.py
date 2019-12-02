from django.shortcuts import get_object_or_404, render, redirect
from django.forms import modelform_factory
from django.http import HttpResponse
import json
import os

class DefaultViews:
    def __init__(self, full_name, app, name, Model, custom=None):
        self.full_name = full_name
        self.title = name
        self.app = app
        self.Model = Model
        self.custom = custom
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
