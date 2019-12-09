from django.shortcuts import render

class SubappViews:
    def __init__(self, app, name, title, models):
        self.app = app
        self.subapp_url ='{app}-{subapp}'.format(app=app, subapp=name)
        self.name = name
        self.title = title
        self.models = models

    def subapp(self, request):
        class ModelInfo:
            def __init__(self, app, subapp, model):
                self.url = '{app}-{subapp}-{model}-index'.format(app=app, subapp=subapp, model=model)
                self.name = model

        modelinfo_list = map(lambda model: ModelInfo(self.app, self.name, model), self.models)

        template = 'subapp.html'
        context = {
            'title': self.title,
            'modelinfo_list': modelinfo_list
        }

        return render(request, template, context)
