from django.urls import path
from .views import Views
from .models import Grid

subapps = [
    { 'full_name': 'Grid Config', 'app': 'config', 'name': 'grid', 'Model': Grid },
]
urlpatterns = []
for subapp in subapps:
    views = Views(**subapp)
    urlpatterns.append(
        path(
            '{name}/'.format(name=subapp['name']),
            views.index,
            name=views.index_url
        )
    )
    urlpatterns.append(
        path(
            '{name}/<int:model_id>/'.format(name=subapp['name']),
            views.detail,
            name=views.detail_url
        )
    )
    urlpatterns.append(
        path(
            '{name}/new/'.format(name=subapp['name']),
            views.new,
            name=views.new_url
        )
    )
