export const ENTER_KEYCODE = 13;
export const ENTER_KEYSTR = 'Enter';

export const NODE_DEFAULT_STR = 'default';
export const NODE_DEFAULT_COLOR = '#2575E2';
export const NODE_ICON_FONT = 'FontAwesome';
export const NODE_ICON_FONT_SIZE = 40;
export const NODE_LABEL_MAX_LENGTH = 40;
export const NODE_LABEL_SPLIT_CHAR = '-';

//Node icons to be displayed in ased on Node types used in app
export const NodeIconMap = new Map([
  //If any of the unicode characters are changed, please also update the
  //comment indicating which FontAwesome icon name is being used.
  //Determine mapping at FA 4.7 docs - https://fontawesome.com/v4.7.0/icons/
  ['default', '\uf111'], //fa-circle
  ['Cluster', '\uf0c2'], //fa-cloud
  ['Namespace', '\uf1c9'], //fa-file-code-o
  ['Service', '\uf085'], //fa-gears
  ['Deployment', '\uf135'], //fa-rocket
  ['ReplicaSet', '\uf1b3'], //fa-cubes
  ['Pod', '\uf1b2'], //fa-cube
  ['Container', '\uf1b2'], //fa-cube
  ['PersistentVolume', '\uf1c0'], //fa-database
  ['PersistentVolumeClaim', '\uf044'], //fa-pencil-square-o
  ['StatefulSet', '\uf0c5'], //fa-copy
  ['Ingress', '\uf090'], //fa-sign-in
  ['Node', '\uf233'], //fa-server
]);

export const NodeStatusColorMap = new Map([
  ['default', NODE_DEFAULT_COLOR], //#e27125 - this orange is the complement to the default blue
  //positive, successful statuses - green
  ['Running', '#25e293'],
  ['Bound', '#25e293'],
  //concerning, risky statuses - yellow
  ['Terminating', '#e2cf25'], //this color doesn't matter since it will be overwritten thru pulsing... but probably good for the legend
  //error, fault statuses - red
  ['Stopped', '#e2254e'],
]);

export const NodeStatusPulseColors = new Map([
  //use RGB array formatted colors in order to be able to easily mix for pulsing
  ['Terminating', [[255,235,87], [176,162,60]]]
]);

export const EdgeColorMap = new Map([
  ['default', NODE_DEFAULT_COLOR],
]);

export const EdgeLabels = ['belongs_to', 'binds', 'claims', 'contains',
  'controlled_by', 'has', 'is_bound_to', 'routes_traffic', 'runs',
  //relationship types used after dgraph metadata changes
  'cluster', 'namespace', 'owner', 'nodename', '~cluster', '~owner', '~namespace'];
