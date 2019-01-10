//TODO:DM - looks like the validate module isn't currently being used in the app, therefore these tests shouldn't really count UNLESS we plan to use it again
import { validateIPaddress } from './validate';

it('should correctly recognize an IP addr', () => {
  expect(validateIPaddress('127.0.0.1')).toBe(true);
});

it('should correctly recognize a non-IP addr', () => {
  expect(validateIPaddress('localhost')).toBe(false);
});