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
import IconButton from '@material-ui/core/IconButton';
import Tooltip from '@material-ui/core/Tooltip';

import { ResultPaginationActionsWrapped } from './ResultPaginationActions';
import { getQueryLayout, rowCellsFromLayout } from './layoutComposer';
import { addResourceIdFilterQSL } from "../../utils/validate";

// Customized table cell theme
export const CustomTableCell = withStyles(theme => ({
  head: {
    backgroundColor: theme.palette.primary.dark,
    color: theme.palette.common.white,
    textTransform: 'capitalize',
    fontSize: 14
  }
}))(TableCell);

const styles = () => ({
  root: {
    width: '100%',
    overflowX: 'auto'
  },
  table: {
    minWidth: 700
  },
  row: {
    height: 36,
    whiteSpace: 'nowrap',
    overflow: 'hidden'
  },
  button: {
    fontFamily: 'FontAwesome'
  },
  link: {
    textDecoration: 'none'
  }
});

class ResultList extends Component {
  handleChangePage = (event, page) => {
    const { submitQuery, query, rowsPerPage } = this.props;
    submitQuery(query.current, page, rowsPerPage);
  };

  handleChangeRowsPerPage = event => {
    const { submitQuery, query, page } = this.props;
    submitQuery(query.current, page, event.target.value);
  };

  renderTableHeader = () => {
    const { query } = this.props;

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

      let columns = [];
      columns.push(
        //header column for set of action icon buttons
        <CustomTableCell key='action-header'>
          <strong>Actions</strong>
        </CustomTableCell>
      );
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

    return tableHeader;
  };

  renderTableRows = () => {
    const { classes, query, onRowClick, selectedIdx } = this.props;
    let tableRows = (
      <TableRow>
        <TableCell />
        <TableCell>No data</TableCell>
        <TableCell />
      </TableRow>
    );

    if (query.results.length > 0) {
      if (query.isQSL) {
        const layout = getQueryLayout(query.current, query.metadata);

        tableRows = query.results.map((item, idx) => {
          let cells = rowCellsFromLayout(item, layout);
          cells.unshift(
            <CustomTableCell key={item.uid}>
              <a className={classes.link} href={'/graph/' + addResourceIdFilterQSL(query.current, item.resourceid)}>
                <Tooltip title="Response Object in Graph View" aria-label="Response in Graph">
                  <IconButton className={classes.button} variant='contained' color='primary' size='small'>
                    {'\uf0e8'}
                  </IconButton>
                </Tooltip>
              </a>
            </CustomTableCell>
          );

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
        tableRows = query.results.map((item, idx) => {
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
    return tableRows;
  };

  render() {
    const { classes, query } = this.props;

    return (
      <Paper className={classes.root} square={true}>
        <Table padding='dense' className={classes.table}>
          <TableHead>{this.renderTableHeader()}</TableHead>
          <TableBody>{this.renderTableRows()}</TableBody>
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
