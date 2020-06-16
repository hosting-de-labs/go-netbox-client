from dcim.models import Site, Rack, DeviceRole, DeviceType, Device, Platform
from ipam.models import IPAddress
from startup_script_utils import load_yaml
import sys

devices = load_yaml('/opt/netbox/initializers/dcim_devices.yml')

if devices is None:
  sys.exit()

handled_attrs = [
  'primary_ip4_id',
  'primary_ip6_id'
]

for params in devices:
  update = False
  new_params = {}
  for field in handled_attrs:
    if field in params:
      update = True
      new_params[field] = params[field]

  if len(new_params) == 0:
    continue

  if update:
    Device.objects.filter(name=params['name']).update(**new_params)
    print("🖥️  Updated device", params['name'])
