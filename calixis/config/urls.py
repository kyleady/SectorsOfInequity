from django.urls import path
from .views import Views
from .models import Grid, Sector, SectorSystem, SectorRoute

subapps = [
    { 'full_name': 'Grid Config',   'app': 'config', 'name': 'grid',   'Model': Grid },
    { 'full_name': 'Sector Config', 'app': 'config', 'name': 'sector', 'Model': Sector, 'SubModels': [SectorSystem, SectorRoute] },
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
    if subapp["name"] == "sector":
        urlpatterns.append(
            path(
                '{name}/<int:model_id>/'.format(name=subapp['name']),
                views.sector_detail,
                name=views.detail_url
            )
        )
    else:
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
    urlpatterns.append(
        path(
            '{name}/test/<int:model_id>/'.format(name=subapp['name']),
            views.test,
            name=views.test_url
        )
    )
