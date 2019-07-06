from django.http import HttpResponse
from inequity.skeleton.SectorSkeleton import SectorSkeleton
from inequity.config.SkeletonConfig import SkeletonConfig

def index(request):
    config = SkeletonConfig()
    sector = SectorSkeleton(config)
    sector.randomize()
    return HttpResponse(sector.debug_display(newline="<br>"))
