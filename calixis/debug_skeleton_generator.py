import logging
from inequity.config.SkeletonConfig import SkeletonConfig
from inequity.skeleton.SectorSkeleton import SectorSkeleton

log = logging.getLogger(__name__)
if __name__ == "__main__":
    config = SkeletonConfig()
    sector = SectorSkeleton(config)
    sector.randomize()
    print(sector.debug_display(newline="\n"))
