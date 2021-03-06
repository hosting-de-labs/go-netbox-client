from ipam.models import RIR
from startup_script_utils import load_yaml
import sys

rirs = load_yaml('/opt/netbox/initializers/ipam_rirs.yml')

if rirs is None:
  sys.exit()

for params in rirs:
  rir, created = RIR.objects.get_or_create(**params)

  if created:
    print("🗺️ Created RIR", rir.name)
