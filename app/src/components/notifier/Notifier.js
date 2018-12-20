import React from 'react';
import {connect} from 'react-redux';
import Snackbar from '@material-ui/core/Snackbar';
import { withStyles } from '@material-ui/core/styles';

const styles = {
    root: {
        background: 'orange'
    }
};

const AutoHideDuration = 10000;

class Notifier extends React.Component {

  state = {
     open: false,
  };

  componentDidMount() {
  }

  componentWillReceiveProps(nextProps) {
    this.setState({
      open:(+new Date() - nextProps.notify.timestamp) < AutoHideDuration,
    });
  }

  handleSnackbarClose() {
    this.setState({
      open:false,
    });
  }

  render() {
    const { classes } = this.props;
    const message = (
      <span
        id="snackbar-message-id">
        {this.props.notify.msg}
      </span>
    );

    //Need to bind the handleSnackbarClose to this class, because it is called from an async timer fn.
    return (
      <Snackbar
        anchorOrigin={{ vertical: 'top', horizontal: 'right' }}
        message={message}
        autoHideDuration={AutoHideDuration}
        onClose={this.handleSnackbarClose.bind(this)}
        open={this.state.open}
        ContentProps={{
            classes: {
                root: classes.root
            },
            'aria-describedby': 'snackbar-message-id',
        }}
      />
    );
  }
}

const mapStoreToProps = store => ({notify: store.notify});

const mapDispatchToProps = dispatch => ({
});

export default connect(
  mapStoreToProps,
  mapDispatchToProps
)(
  withStyles(styles)(Notifier)
);
