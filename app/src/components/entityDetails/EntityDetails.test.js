import React from 'react';
import { render, unmountComponentAtNode } from 'react-dom';
import { shallow } from 'enzyme';

import EntityDetails from './EntityDetails';

const div = document.createElement('div');

it('shallow renders entity details component', () => {
  shallow(<EntityDetails selectedObj={{}}/>);
});

it('deep renders entity details component', () => {
  render(<EntityDetails selectedObj={{}}/>, div);
  unmountComponentAtNode(div);
});