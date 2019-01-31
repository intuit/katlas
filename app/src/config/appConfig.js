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
  ['cluster', '\uf0c2'], //fa-cloud
  ['namespace', '\uf1c9'], //fa-file-code-o
  ['service', '\uf085'], //fa-gears
  ['deployment', '\uf135'], //fa-rocket
  ['replicaset', '\uf1b3'], //fa-cubes
  ['pod', '\uf1b2'], //fa-cube
  ['container', '\uf1b2'], //fa-cube
  ['persistentvolume', '\uf1c0'], //fa-database
  ['persistentvolumeclaim', '\uf044'], //fa-pencil-square-o
  ['statefulset', '\uf0c5'], //fa-copy
  ['ingress', '\uf090'], //fa-sign-in
  ['node', '\uf233'], //fa-server
  ['daemonset', '\uf2ac'], //fa-snapchat-ghost
  ['application', '\uf0e4'] //fa-tachometer
]);

export const NodeStatusColorMap = new Map([
  ['default', NODE_DEFAULT_COLOR], //#e27125 - this orange is the complement to the default blue
  //positive, successful statuses - green
  ['Running', '#25e293'],
  ['Bound', '#25e293'],
  //concerning, risky statuses - yellow
  ['Terminating', '#e2cf25'], //this color will be overwritten thru pulsing, but is still needed for the legend
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
  //relationship types used after dgraph metadata changes, not currently using
  //'~cluster', '~namesapce' since they often have hundreds of children
  'cluster', 'namespace', 'owner', 'nodename', '~owner', '~application'];
