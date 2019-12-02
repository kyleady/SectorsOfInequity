from django.urls import path
from .views.default import DefaultViews
from .views.sector import SectorViews
from .models import Grid, Sector, SectorSystem, SectorRoute, Region, WeightedRegion

subapps = [
    { 'full_name': 'Region Config', 'app': 'config', 'name': 'region', 'Model': Region },
    { 'full_name': 'Weighted Region', 'app': 'config', 'name': 'weighted-region', 'Model': WeightedRegion },
    { 'full_name': 'Grid Config',   'app': 'config', 'name': 'grid',   'Model': Grid },
    { 'full_name': 'Sector Config', 'app': 'config', 'name': 'sector', 'Model': Sector, 'custom': { 'SubModels': [SectorSystem, SectorRoute], 'Grid': Grid }, 'Views': SectorViews },
]
urlpatterns = []
for subapp in subapps:
    if 'Views' in subapp:
        CustomViews = subapp['Views']
        del subapp['Views']
        views = CustomViews(**subapp)
    else:
        views = DefaultViews(**subapp)

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
    urlpatterns.append(
        path(
            '{name}/delete/<int:model_id>'.format(name=subapp['name']),
            views.delete,
            name=views.delete_url
        )
    )
