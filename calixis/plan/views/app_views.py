from django.shortcuts import render

class AppViews:
    def __init__(self, name, title, subapps):
        self.subapps = subapps
        self.app_url = name
        self.name = name
        self.title = title
        self.subapps = subapps

    def app(self, request):
        class SubappInfo:
            def __init__(self, app, subapp):
                self.url = '{app}-{subapp}'.format(app=app, subapp=subapp)
                self.name = subapp

        subappinfo_list = map(lambda subapp: SubappInfo(self.name, subapp), self.subapps)

        template = 'app.html'
        context = {
            'title': self.title,
            'subappinfo_list': subappinfo_list
        }

        return render(request, template, context)
