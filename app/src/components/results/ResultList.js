import React, { Component } from 'react';
import { withStyles } from '@material-ui/core/styles';
import { Link } from 'react-router-dom';
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableHead from '@material-ui/core/TableHead';
import TableRow from '@material-ui/core/TableRow';
import Paper from '@material-ui/core/Paper';

// Customized table cell theme
const CustomTableCell = withStyles(theme => ({
  head: {
    backgroundColor: theme.palette.primary.dark,
    color: theme.palette.common.white,
    fontSize: 14,
  },
}))(TableCell);

const styles = theme => ({
  root: {
    width: '100%',
    overflowX: 'auto',
  },
  table: {
    minWidth: 700,
  },
  row: {
    height: 36,
  }
});

class ResultList extends Component {

  render() {
    const { classes, query, onRowClick, selectedIdx } = this.props;

    return (
      <Paper className={classes.root} square={true}>
        <Table padding='dense' className={classes.table}>
          <TableHead>
            <TableRow>
              <CustomTableCell>Type</CustomTableCell>
              <CustomTableCell>Name</CustomTableCell>
              <CustomTableCell>Namespace</CustomTableCell>
              <CustomTableCell>Created At</CustomTableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {(query.results.length > 0) ? (
              query.results.map((item, idx) => {
                return (
                  <TableRow hover key={item.uid} className={classes.row}
                    onClick={event => onRowClick(event, idx)} selected={selectedIdx === idx}>
                    <CustomTableCell component='th' scope='row'>
                      {item.objtype}
                    </CustomTableCell>
                    <CustomTableCell>
                      <Link
                        to={{
                          pathname: '/graph/' + item.uid,
                          state: { selectedObj: query.results[selectedIdx] }
                        }}>
                        {item.name}
                      </Link>
                    </CustomTableCell>
                    <CustomTableCell>{item.namespace ? item.namespace[0].name : ''}</CustomTableCell>
                    <CustomTableCell>{item.creationtime}</CustomTableCell>
                  </TableRow>
                );
              })
            ) : ( //TODO:DM determine if there is a more elegant 'toggle' pattern suggested in React/jsx community
                <TableRow>
                  <TableCell />
                  <TableCell>No data</TableCell>
                  <TableCell />
                </TableRow>
              )}
          </TableBody>
        </Table>
      </Paper>
    )
  }
}

export default withStyles(styles)(ResultList);