from django.urls import path
from .views.default import DefaultViews
from .views.asset import AssetViews
from .views.asset_grid import AssetGridViews
from .views.app_views import AppViews
from .views.subapp_views import SubappViews
from .models import *

app = 'plan'
subapp = 'asset'
#  'custom': { 'Config': Grid_Sector, 'Job': Job, 'detail_template': 'asset_sector_detail.html', 'no_form': True }, 'Views': AssetSectorViews
asset_models = [
    { 'full_name': 'Asset Group', 'app': app, 'subapp': subapp, 'name': 'asset-group', 'Model': Asset_Group },
    { 'full_name': 'Asset', 'app': app, 'subapp': subapp, 'name': 'asset', 'Model': Asset,
      'custom': {
        'Perterbation_Model': Perterbation, 'Config_Name_Model': Config_Name, 'Job': Job, 'detail_template': 'asset_detail.html', 'no_form': True
      },
      'Views': AssetViews
    },
    { 'full_name': 'Asset Grid', 'app': app, 'subapp': subapp, 'name': 'asset-grid', 'Model': Asset_Grid,
        'custom': {
           'SubModels': [Asset_Node, Asset_Connection], 'Grid': Asset_Grid, 'Job': Job, 'detail_template': 'asset_grid_detail.html', 'no_form': True
        },
        'Views': AssetGridViews
    },
    { 'full_name': 'Asset Node', 'app': app, 'subapp': subapp, 'name': 'asset-node', 'Model': Asset_Node },
    { 'full_name': 'Asset Connection', 'app': app, 'subapp': subapp, 'name': 'asset-connection', 'Model': Asset_Connection },
]

subapp = 'config'
config_models = [
    { 'full_name': 'Tag', 'app': app, 'subapp': subapp, 'name': 'tag', 'Model': Tag },
    { 'full_name': 'Config Name', 'app': app, 'subapp': subapp, 'name': 'config-name', 'Model': Config_Name },
    { 'full_name': 'Config Asset', 'app': app, 'subapp': subapp, 'name': 'config-asset', 'Model': Config_Asset },
    { 'full_name': 'Config Group',  'app': app, 'subapp': subapp, 'name': 'config-sub-asset',  'Model': Config_Group },
    { 'full_name': 'Config Grid', 'app': app, 'subapp': subapp, 'name': 'config-grid', 'Model': Config_Grid },
]

subapp = 'detail'
detail_models = [
    { 'full_name': 'Detail', 'app': app, 'subapp': subapp, 'name': 'detail', 'Model': Detail,
        'custom': { 'detail_template': 'detail_detail_detail.html', 'no_form': True } },
]

subapp = 'inspiration'
inspiration_models = [
    { 'full_name': 'Inspiration', 'app': app, 'subapp': subapp, 'name': 'inspiration', 'Model': Inspiration },
    { 'full_name': 'Inspiration Table', 'app': app, 'subapp': subapp, 'name': 'inspiration-table', 'Model': Inspiration_Table },
    { 'full_name': 'Inspiration Extra', 'app': app, 'subapp': subapp, 'name': 'inspiration-extra', 'Model': Inspiration_Extra },
    { 'full_name': 'Weighted Inspiration', 'app': app, 'subapp': subapp, 'name': 'weighted-inspiration', 'Model': Weighted_Inspiration },
]


subapp = 'roll'
roll_models = [
    { 'full_name': 'Roll', 'app': app, 'subapp': subapp, 'name': 'roll', 'Model': Roll },
]

subapp = 'job'
job_models = [
    { 'full_name': 'Job', 'app': app, 'subapp': subapp, 'name': 'job', 'Model': Job },
]

all_models = (
    asset_models +
    config_models +
    detail_models +
    inspiration_models +
    roll_models +
    job_models)

urlpatterns = []
for model in all_models:
    if 'Views' in model:
        CustomViews = model['Views']
        del model['Views']
        views = CustomViews(**model)
    else:
        views = DefaultViews(**model)

    urlpatterns.append(
        path(
            '{subapp}/{name}/'.format(
                subapp=model['subapp'],
                name=model['name'],
            ),
            views.index,
            name=views.index_url
        )
    )
    urlpatterns.append(
        path(
            '{subapp}/{name}/<int:model_id>/'.format(
                subapp=model['subapp'],
                name=model['name'],
            ),
            views.detail,
            name=views.detail_url
        )
    )
    urlpatterns.append(
        path(
            '{subapp}/{name}/new/'.format(
                subapp=model['subapp'],
                name=model['name'],
            ),
            views.new,
            name=views.new_url
        )
    )
    urlpatterns.append(
        path(
            '{subapp}/{name}/delete/<int:model_id>'.format(
                subapp=model['subapp'],
                name=model['name'],
            ),
            views.delete,
            name=views.delete_url
        )
    )

app_views = AppViews(title='Sector Planning', name='plan', subapps=[
                'asset',
                'config',
                'detail',
                'inspiration',
                'roll',
                'job'
            ])

urlpatterns.append(
    path(
        '',
        app_views.app,
        name=app_views.app_url
    )
)

all_subapp_views = [
                SubappViews(title='Assets', name=asset_models[0]['subapp'], app=asset_models[0]['app'], models=[
                    model['name'] for model in asset_models
                ]),
                SubappViews(title='Config',name=config_models[0]['subapp'], app=config_models[0]['app'], models=[
                    model['name'] for model in config_models
                ]),
                SubappViews(title='Details', name=detail_models[0]['subapp'], app=detail_models[0]['app'], models=[
                    model['name'] for model in detail_models
                ]),
                SubappViews(title='Inspiration', name=inspiration_models[0]['subapp'], app=inspiration_models[0]['app'], models=[
                    model['name'] for model in inspiration_models
                ]),
                SubappViews(title='Rolls', name=roll_models[0]['subapp'], app=roll_models[0]['app'], models=[
                    model['name'] for model in roll_models
                ]),
                SubappViews(title='Job', name=job_models[0]['subapp'], app=job_models[0]['app'], models=[
                    model['name'] for model in job_models
                ]),
]

for subapp_views in all_subapp_views:
    urlpatterns.append(
        path(
            '{subapp}'.format(subapp=subapp_views.name),
            subapp_views.subapp,
            name=subapp_views.subapp_url
        )
    )
