import { applyMiddleware, createStore, compose } from 'redux';
import rootReducer from '../reducers/rootReducer';
import thunk from 'redux-thunk';

//allow debug use of Chrome DevTool for Redux apps
const composeEnhancers = window.__REDUX_DEVTOOLS_EXTENSION_COMPOSE__ || compose;

export default function configureStore() {
  return createStore(
    rootReducer, composeEnhancers(
      applyMiddleware(thunk)
    )
  );
}