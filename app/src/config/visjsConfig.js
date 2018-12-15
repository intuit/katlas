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
        improvedLayout: true,
    },
    physics: {
        enabled: false,
    }
};