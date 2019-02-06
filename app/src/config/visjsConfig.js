//TODO:DM - come up with better name for this file/object; and move to be closer to graph component
export const options = {
  nodes: {
    shape: 'icon',
    scaling: {
      max: 20,
      min: 20,
      label: {
        enabled: true,
        min: 14,
        max: 14
      }
    },
    font: {
      size: 14
    },
    margin: {
      top: 25
    }
  },
  height: "100%",
  width: "100%",
  interaction: {
    hover: true,
    keyboard: {
      enabled: true,
      bindToWindow: false
    },
    navigationButtons: true,
    tooltipDelay: 1000000,
    zoomView: true
  },
  layout: {
    randomSeed: 42,
    improvedLayout: true,
    hierarchical: {
      enabled: true,
      nodeSpacing: 150,
      blockShifting: true,
      edgeMinimization: true,
      sortMethod: 'hubsize',
      direction: 'UD'
    }
  },
  physics: {
    enabled: false
  }
};