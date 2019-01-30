import React, { Component } from 'react';
import { withStyles } from '@material-ui/core/styles';
import { Link } from 'react-router-dom';
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableHead from '@material-ui/core/TableHead';
import TableRow from '@material-ui/core/TableRow';
import TableFooter from '@material-ui/core/TableFooter';
import Paper from '@material-ui/core/Paper';
import TablePagination from '@material-ui/core/TablePagination';

import { ResultPaginationActionsWrapped } from './ResultPaginationActions';
import { getQueryLayout, rowCellsFromLayout } from './layoutComposer';

// Customized table cell theme
export const CustomTableCell = withStyles(theme => ({
  head: {
    backgroundColor: theme.palette.primary.dark,
    color: theme.palette.common.white,
    textTransform: 'capitalize',
    fontSize: 14
  }
}))(TableCell);

const styles = theme => ({
  root: {
    width: '100%',
    overflowX: 'auto'
  },
  table: {
    minWidth: 700
  },
  row: {
    height: 36
  }
});

class ResultList extends Component {
  handleChangePage = (event, page) => {
    const {
      updatePagination,
      query: { rowsPerPage }
    } = this.props;
    updatePagination(page, rowsPerPage);
  };

  handleChangeRowsPerPage = event => {
    const {
      updatePagination,
      query: { page }
    } = this.props;
    updatePagination(page, event.target.value);
  };

  render() {
    const { classes, query, onRowClick, selectedIdx } = this.props;
    let tableHeader = (
      <TableRow>
        <CustomTableCell>Type</CustomTableCell>
        <CustomTableCell>Name</CustomTableCell>
        <CustomTableCell>Namespace</CustomTableCell>
        <CustomTableCell>Created At</CustomTableCell>
      </TableRow>
    );

    if (query.isQSL) {
      const layout = getQueryLayout(query.current, query.metadata);

      const columns = [];
      for (let objType in layout) {
        const fields = layout[objType];
        for (let fieldname in fields) {
          const representer = fields[fieldname];
          columns.push(
            <CustomTableCell key={`${objType}-${fieldname}`}>
              <strong>{objType}</strong>
              <br />
              {representer.displayName}
            </CustomTableCell>
          );
        }
      }
      tableHeader = <TableRow>{columns}</TableRow>;
    }

    let tableRow = (
      <TableRow>
        <TableCell />
        <TableCell>No data</TableCell>
        <TableCell />
      </TableRow>
    );

    if (query.results.length > 0) {
      if (query.isQSL) {
        const layout = getQueryLayout(query.current, query.metadata);

        tableRow = query.results.map((item, idx) => {
          const cells = rowCellsFromLayout(item, layout);

          return (
            <TableRow
              hover
              key={item.uid}
              className={classes.row}
              onClick={event => onRowClick(event, idx)}
              selected={selectedIdx === idx}
            >
              {cells}
            </TableRow>
          );
        });
      } else {
        tableRow = query.results.map((item, idx) => {
          return (
            <TableRow
              hover
              key={item.uid}
              className={classes.row}
              onClick={event => onRowClick(event, idx)}
              selected={selectedIdx === idx}
            >
              <CustomTableCell component='th' scope='row'>
                {item.objtype}
              </CustomTableCell>
              <CustomTableCell>
                <Link
                  to={{
                    pathname: '/graph/' + item.uid,
                    state: { selectedObj: query.results[selectedIdx] }
                  }}
                >
                  {item.name}
                </Link>
              </CustomTableCell>
              <CustomTableCell>
                {item.namespace ? item.namespace[0].name : ''}
              </CustomTableCell>
              <CustomTableCell>{item.creationtime}</CustomTableCell>
            </TableRow>
          );
        });
      }
    }

    return (
      <Paper className={classes.root} square={true}>
        <Table padding='dense' className={classes.table}>
          <TableHead>{tableHeader}</TableHead>
          <TableBody>{tableRow}</TableBody>
          <TableFooter>
            <TableRow>
              <TablePagination
                rowsPerPageOptions={[10, 25, 50]}
                count={query.count}
                rowsPerPage={query.rowsPerPage}
                page={query.page}
                onChangePage={this.handleChangePage}
                onChangeRowsPerPage={this.handleChangeRowsPerPage}
                ActionsComponent={ResultPaginationActionsWrapped}
              />
            </TableRow>
          </TableFooter>
        </Table>
      </Paper>
    );
  }
}

export default withStyles(styles)(ResultList);
