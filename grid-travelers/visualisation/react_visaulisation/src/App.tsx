import Grid from './Grid.tsx';
import './App.css';

function App() {
    return (
        <div className="grid-container">
            <Grid rows={10} cols={15} />
        </div>
    );
}

export default App;