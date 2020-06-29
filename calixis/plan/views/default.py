from django.shortcuts import get_object_or_404, render, redirect
from django.core.paginator import Paginator
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
        self.custom = custom or {}
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
        list_of_objects = self.Model.objects.all().order_by('id')
        paginator = Paginator(list_of_objects, 25)
        page_number = request.GET.get('page')
        page_obj = paginator.get_page(page_number)

        if self.custom.get('detail_template', False):
            template = self.custom['index_template']
        else:
            template = 'index.html'
        context = {
            'new_url': self.new_url,
            'detail_url': self.detail_url,
            'full_name': self.full_name,
            'page_obj': page_obj
        }

        return render(request, template, context)

    def detail(self, request, model_id):
        if self.custom.get('detail_template', False):
            template = self.custom['detail_template']
        else:
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

        if self.custom.get('no_form', False):
            del context['form']

        return render(request, template, context)

    def new(self, request):
        if self.custom.get('detail_template', False):
            template = self.custom['new_template']
        else:
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
