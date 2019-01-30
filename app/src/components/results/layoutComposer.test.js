import { getQueryLayout, rowCellsFromLayout } from './layoutComposer';

const queryStr =
  'deployment{*}.replicaset[@count(pod)<3]{*}.pod{@name,@resourceid}';
const metadata = {
  deployment: {
    uid: '0xfb80d',
    name: 'deployment',
    fields: [
      {
        fieldname: 'labels',
        fieldtype: 'json',
        mandatory: false,
        cardinality: 'one'
      },
      {
        fieldname: 'name',
        fieldtype: 'string',
        mandatory: true,
        cardinality: 'one'
      },
      {
        fieldname: 'numreplicas',
        fieldtype: 'int',
        mandatory: true,
        cardinality: 'one'
      },
      {
        fieldname: 'resourceversion',
        fieldtype: 'string',
        mandatory: true,
        cardinality: 'one'
      },
      {
        fieldname: 'strategy',
        fieldtype: 'string',
        mandatory: true,
        cardinality: 'one'
      },
      {
        fieldname: 'creationtime',
        fieldtype: 'string',
        mandatory: true,
        cardinality: 'one'
      },
      {
        fieldname: 'cluster',
        fieldtype: 'relationship',
        mandatory: true,
        refdatatype: 'cluster',
        cardinality: 'one'
      },
      {
        fieldname: 'availablereplicas',
        fieldtype: 'int',
        mandatory: true,
        cardinality: 'one'
      },
      {
        fieldname: 'objtype',
        fieldtype: 'string',
        mandatory: true,
        cardinality: 'one'
      },
      {
        fieldname: 'namespace',
        fieldtype: 'relationship',
        mandatory: true,
        refdatatype: 'namespace',
        cardinality: 'one'
      },
      {
        fieldname: 'resourceid',
        fieldtype: 'string',
        mandatory: false,
        cardinality: 'one'
      },
      {
        fieldname: 'application',
        fieldtype: 'relationship',
        mandatory: false,
        refdatatype: 'application',
        cardinality: 'many'
      },
      {
        fieldname: 'k8sobj',
        fieldtype: 'true',
        mandatory: true,
        cardinality: 'one'
      }
    ]
  },
  replicaset: {
    uid: '0xfb821',
    name: 'replicaset',
    fields: [
      {
        fieldname: 'resourceid',
        fieldtype: 'string',
        mandatory: false,
        cardinality: 'one'
      },
      {
        fieldname: 'cluster',
        fieldtype: 'relationship',
        mandatory: true,
        refdatatype: 'cluster',
        cardinality: 'one'
      },
      {
        fieldname: 'k8sobj',
        fieldtype: 'true',
        mandatory: true,
        cardinality: 'one'
      },
      {
        fieldname: 'name',
        fieldtype: 'string',
        mandatory: true,
        cardinality: 'one'
      },
      {
        fieldname: 'labels',
        fieldtype: 'json',
        mandatory: false,
        cardinality: 'one'
      },
      {
        fieldname: 'resourceversion',
        fieldtype: 'string',
        mandatory: true,
        cardinality: 'one'
      },
      {
        fieldname: 'creationtime',
        fieldtype: 'string',
        mandatory: true,
        cardinality: 'one'
      },
      {
        fieldname: 'objtype',
        fieldtype: 'string',
        mandatory: true,
        cardinality: 'one'
      },
      {
        fieldname: 'namespace',
        fieldtype: 'relationship',
        mandatory: true,
        refdatatype: 'namespace',
        cardinality: 'one'
      },
      {
        fieldname: 'owner',
        fieldtype: 'relationship',
        mandatory: false,
        refdatatype: 'deployment',
        cardinality: 'one'
      },
      {
        fieldname: 'numreplicas',
        fieldtype: 'int',
        mandatory: true,
        cardinality: 'one'
      },
      {
        fieldname: 'podspec',
        fieldtype: 'json',
        mandatory: false,
        cardinality: 'one'
      }
    ]
  },
  pod: {
    uid: '0xfb858',
    name: 'pod',
    fields: [
      {
        fieldname: 'name',
        fieldtype: 'string',
        mandatory: true,
        cardinality: 'one'
      },
      {
        fieldname: 'phase',
        fieldtype: 'string',
        mandatory: false,
        cardinality: 'one'
      },
      {
        fieldname: 'ip',
        fieldtype: 'string',
        mandatory: true,
        cardinality: 'one'
      },
      {
        fieldname: 'containers',
        fieldtype: 'json',
        mandatory: false,
        cardinality: 'many'
      },
      {
        fieldname: 'starttime',
        fieldtype: 'string',
        mandatory: false,
        cardinality: 'one'
      },
      {
        fieldname: 'resourceid',
        fieldtype: 'string',
        mandatory: false,
        cardinality: 'one'
      },
      {
        fieldname: 'cluster',
        fieldtype: 'relationship',
        mandatory: true,
        refdatatype: 'cluster',
        cardinality: 'one'
      },
      {
        fieldname: 'namespace',
        fieldtype: 'relationship',
        mandatory: true,
        refdatatype: 'namespace',
        cardinality: 'one'
      },
      {
        fieldname: 'labels',
        fieldtype: 'json',
        mandatory: false,
        cardinality: 'one'
      },
      {
        fieldname: 'resourceversion',
        fieldtype: 'string',
        mandatory: true,
        cardinality: 'one'
      },
      {
        fieldname: 'nodename',
        fieldtype: 'relationship',
        mandatory: true,
        refdatatype: 'node',
        cardinality: 'one'
      },
      {
        fieldname: 'ownertype',
        fieldtype: 'string',
        mandatory: false,
        cardinality: 'one'
      },
      {
        fieldname: 'creationtime',
        fieldtype: 'string',
        mandatory: true,
        cardinality: 'one'
      },
      {
        fieldname: 'k8sobj',
        fieldtype: 'true',
        mandatory: true,
        cardinality: 'one'
      },
      {
        fieldname: 'objtype',
        fieldtype: 'string',
        mandatory: true,
        cardinality: 'one'
      },
      {
        fieldname: 'volumes',
        fieldtype: 'json',
        mandatory: false,
        cardinality: 'many'
      },
      {
        fieldname: 'owner',
        fieldtype: 'relationship',
        mandatory: false,
        refdatatype: 'replicaset,daemonset,statefulset',
        cardinality: 'one'
      }
    ]
  }
};
const dataItem = {
  availablereplicas: '1',
  creationtime: '2019-01-18T01:18:36Z',
  k8sobj: 'k8sobj',
  labels:
    '{"app":"foremast-service","applications.argoproj.io/app-name":"foremast-brain-usw2-ppd"}',
  name: 'foremast-service',
  numreplicas: '1',
  objtype: 'deployment',
  resourceid:
    'deployment:paas-preprod-west2.cluster.k8s.local:dev-containers-foremast-brain-usw2-ppd-ppd:foremast-service',
  resourceversion: '90347576',
  strategy: 'RollingUpdate',
  uid: '0x10f0f9',
  '~owner': [
    {
      creationtime: '2019-01-18T17:41:01Z',
      k8sobj: 'k8sobj',
      labels:
        '{"app":"foremast-service","applications.argoproj.io/app-name":"foremast-brain-usw2-ppd","pod-template-hash":"3387190116"}',
      name: 'foremast-service-77dc5f455b',
      numreplicas: '1',
      objtype: 'replicaset',
      resourceid:
        'replicaset:paas-preprod-west2.cluster.k8s.local:dev-containers-foremast-brain-usw2-ppd-ppd:foremast-service-77dc5f455b',
      resourceversion: '90347574',
      uid: '0x120272',
      '~owner': [
        {
          containers:
            '[{"name":"service","image":"docker.artifactory.a.intuit.com/foremast/foremast-service:0.0.5","ports":[{"name":"http","containerPort":8099,"protocol":"TCP"}],"env":[{"name":"ELASTIC_URL","value":"http://elasticsearch-discovery.dev-containers-foremast-brain-usw2-ppd-ppd.svc.cluster.local:9200/"}],"resources":{"limits":{"cpu":"100m","memory":"30Mi"},"requests":{"cpu":"100m","memory":"20Mi"}},"volumeMounts":[{"name":"default-token-zxbmj","readOnly":true,"mountPath":"/var/run/secrets/kubernetes.io/serviceaccount"}],"terminationMessagePath":"/dev/termination-log","terminationMessagePolicy":"File","imagePullPolicy":"Always"}]',
          creationtime: '2019-01-20T08:35:09Z',
          ip: '100.100.135.8',
          k8sobj: 'k8sobj',
          labels:
            '{"app":"foremast-service","applications.argoproj.io/app-name":"foremast-brain-usw2-ppd","pod-template-hash":"3387190116"}',
          name: 'foremast-service-77dc5f455b-5lpsg',
          nodename: [
            {
              k8sobj: 'k8sobj',
              name: 'ip-10-83-122-52.us-west-2.compute.internal',
              objtype: 'node',
              resourceid:
                'node:paas-preprod-west2.cluster.k8s.local:ip-10-83-122-52.us-west-2.compute.internal',
              resourceversion: '0',
              uid: '0x10c9ed'
            }
          ],
          objtype: 'pod',
          ownertype: 'replicaset',
          phase: 'Running',
          resourceid:
            'pod:paas-preprod-west2.cluster.k8s.local:dev-containers-foremast-brain-usw2-ppd-ppd:foremast-service-77dc5f455b-5lpsg',
          resourceversion: '90347573',
          starttime: '2019-01-20T08:35:09Z',
          uid: '0x18b916',
          volumes:
            '[{"name":"default-token-zxbmj","secret":{"secretName":"default-token-zxbmj","defaultMode":420}}]'
        }
      ]
    },
    {
      creationtime: '2019-01-18T01:18:36Z',
      k8sobj: 'k8sobj',
      labels: '{"app":"foremast-service","pod-template-hash":"1964667380"}',
      name: 'foremast-service-5fb8bbc7d4',
      numreplicas: '0',
      objtype: 'replicaset',
      resourceid:
        'replicaset:paas-preprod-west2.cluster.k8s.local:dev-containers-foremast-brain-usw2-ppd-ppd:foremast-service-5fb8bbc7d4',
      resourceversion: '88371578',
      uid: '0x17f5ec'
    }
  ]
};

describe('layout composer', () => {
  it('can compose the result table layout', () => {
    const layout = getQueryLayout(queryStr, metadata);
    // console.log(layout);
  });

  it('can render data row based on layout', () => {
    const layout = getQueryLayout(queryStr, metadata);
    const row = rowCellsFromLayout(dataItem, layout);
    console.log(row);
  });
});
