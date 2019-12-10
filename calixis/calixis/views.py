from django.shortcuts import render

class HomeViews:
    def __init__(self, apps):
        self.apps = apps
        self.home_url = 'home'
        self.title = 'Sectors of Inequity'

    def home(self, request):
        template = 'home.html'
        context = {
            'title': self.title,
            'app_list': self.apps
        }

        return render(request, template, context)
