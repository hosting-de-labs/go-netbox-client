from dcim.models import DeviceType, Manufacturer, Region
from tenancy.models import Tenant
from extras.models import CustomField, CustomFieldValue
from startup_script_utils import load_yaml
import sys

device_types = load_yaml('/opt/netbox/initializers/dcim_device_types.yml')

if device_types is None:
  sys.exit()

required_assocs = {
  'manufacturer': (Manufacturer, 'name')
}

optional_assocs = {
  'region': (Region, 'name'),
  'tenant': (Tenant, 'name')
}

for params in device_types:
  custom_fields = params.pop('custom_fields', None)

  for assoc, details in required_assocs.items():
    model, field = details
    query = { field: params.pop(assoc) }

    params[assoc] = model.objects.get(**query)

  for assoc, details in optional_assocs.items():
    if assoc in params:
      model, field = details
      query = { field: params.pop(assoc) }

      params[assoc] = model.objects.get(**query)

  device_type, created = DeviceType.objects.get_or_create(**params)

  if created:
    if custom_fields is not None:
      for cf_name, cf_value in custom_fields.items():
        custom_field = CustomField.objects.get(name=cf_name)
        custom_field_value = CustomFieldValue.objects.create(
          field=custom_field,
          obj=device_type,
          value=cf_value
        )

        device_type.custom_field_values.add(custom_field_value)

    print("🔡 Created device type", device_type.manufacturer, device_type.model)
