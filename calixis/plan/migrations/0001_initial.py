# Generated by Django 2.2.6 on 2019-12-08 06:07

from django.db import migrations, models
import django.db.models.deletion


class Migration(migrations.Migration):

    initial = True

    dependencies = [
    ]

    operations = [
        migrations.CreateModel(
            name='Config_Grid',
            fields=[
                ('id', models.AutoField(auto_created=True, primary_key=True, serialize=False, verbose_name='ID')),
                ('name', models.CharField(default='-', max_length=25)),
                ('height', models.PositiveSmallIntegerField(blank=True, default=20)),
                ('width', models.PositiveSmallIntegerField(blank=True, default=20)),
                ('connectionRange', models.PositiveSmallIntegerField(blank=True, default=5)),
                ('populationRate', models.FloatField(blank=True, default=0.5)),
                ('connectionRate', models.FloatField(blank=True, default=0.5)),
                ('rangeRateMultiplier', models.FloatField(blank=True, default=0.5)),
                ('smoothingFactor', models.FloatField(blank=True, default=0.5)),
            ],
            options={
                'abstract': False,
            },
        ),
        migrations.CreateModel(
            name='Config_Region',
            fields=[
                ('id', models.AutoField(auto_created=True, primary_key=True, serialize=False, verbose_name='ID')),
                ('name', models.CharField(default='-', max_length=25)),
            ],
            options={
                'abstract': False,
            },
        ),
        migrations.CreateModel(
            name='Config_Sector',
            fields=[
                ('id', models.AutoField(auto_created=True, primary_key=True, serialize=False, verbose_name='ID')),
                ('name', models.CharField(default='-', max_length=39)),
            ],
            options={
                'abstract': False,
            },
        ),
        migrations.CreateModel(
            name='Config_System',
            fields=[
                ('id', models.AutoField(auto_created=True, primary_key=True, serialize=False, verbose_name='ID')),
                ('name', models.CharField(default='-', max_length=25)),
            ],
            options={
                'abstract': False,
            },
        ),
        migrations.CreateModel(
            name='Inspiration_System',
            fields=[
                ('id', models.AutoField(auto_created=True, primary_key=True, serialize=False, verbose_name='ID')),
                ('name', models.CharField(default='-', max_length=25)),
                ('description', models.CharField(default='-', max_length=1000)),
            ],
            options={
                'abstract': False,
            },
        ),
        migrations.CreateModel(
            name='Perterbation_System',
            fields=[
                ('config_region_ptr', models.OneToOneField(auto_created=True, on_delete=django.db.models.deletion.CASCADE, parent_link=True, primary_key=True, serialize=False, to='plan.Config_Region')),
            ],
            options={
                'abstract': False,
            },
            bases=('plan.config_region',),
        ),
        migrations.CreateModel(
            name='Weighted_Inspiration_Region',
            fields=[
                ('id', models.AutoField(auto_created=True, primary_key=True, serialize=False, verbose_name='ID')),
                ('weight', models.SmallIntegerField()),
                ('parent', models.ForeignKey(on_delete=django.db.models.deletion.CASCADE, to='plan.Config_System')),
                ('value', models.ForeignKey(on_delete=django.db.models.deletion.CASCADE, to='plan.Inspiration_System')),
            ],
            options={
                'abstract': False,
            },
        ),
        migrations.CreateModel(
            name='Weighted_Config_Region',
            fields=[
                ('id', models.AutoField(auto_created=True, primary_key=True, serialize=False, verbose_name='ID')),
                ('weight', models.SmallIntegerField()),
                ('parent', models.ForeignKey(on_delete=django.db.models.deletion.CASCADE, to='plan.Config_Grid')),
                ('value', models.ForeignKey(on_delete=django.db.models.deletion.CASCADE, to='plan.Config_Region')),
            ],
            options={
                'abstract': False,
            },
        ),
        migrations.CreateModel(
            name='Config_Sector_System',
            fields=[
                ('id', models.AutoField(auto_created=True, primary_key=True, serialize=False, verbose_name='ID')),
                ('x', models.PositiveSmallIntegerField()),
                ('y', models.PositiveSmallIntegerField()),
                ('region', models.ForeignKey(on_delete=django.db.models.deletion.CASCADE, to='plan.Config_Region')),
                ('sector', models.ForeignKey(on_delete=django.db.models.deletion.CASCADE, to='plan.Config_Sector')),
            ],
        ),
        migrations.CreateModel(
            name='Config_Sector_Route',
            fields=[
                ('id', models.AutoField(auto_created=True, primary_key=True, serialize=False, verbose_name='ID')),
                ('end', models.ForeignKey(on_delete=django.db.models.deletion.CASCADE, related_name='end', to='plan.Config_Sector_System')),
                ('start', models.ForeignKey(on_delete=django.db.models.deletion.CASCADE, related_name='start', to='plan.Config_Sector_System')),
            ],
        ),
        migrations.AddField(
            model_name='config_region',
            name='system',
            field=models.ForeignKey(on_delete=django.db.models.deletion.CASCADE, to='plan.Config_System'),
        ),
        migrations.AddField(
            model_name='inspiration_system',
            name='perterbation',
            field=models.ForeignKey(on_delete=django.db.models.deletion.CASCADE, to='plan.Perterbation_System'),
        ),
    ]