from django.shortcuts import get_object_or_404, render, redirect
from django.forms import modelform_factory

class Views:
    def __init__(self, full_name, app, name, Model):
        self.full_name = full_name
        self.title = name
        self.app = app
        self.Model = Model
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
        self.Form = modelform_factory(self.Model, fields='__all__')

    def index(self, request):
        template = '{title}/index.html'.format(
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
        template = '{title}/detail.html'.format(
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

    def new(self, request):
        template = '{title}/detail.html'.format(
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
