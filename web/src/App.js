import React, { useState, useEffect } from 'react';
import { BrowserRouter as Router, Routes, Route, Link, useParams, useNavigate } from 'react-router-dom';
import './App.css';

function App() {
  return (
    <Router>
      <div className="App">
        <header className="App-header">
          <h1>Meme Generator</h1>
          <nav>
            <Link to="/">Gallery</Link> | <Link to="/create">Create Meme</Link> | <Link to="/templates">Templates</Link> | <Link to="/templates/create">Create Template</Link>
          </nav>
        </header>
        <main>
          <Routes>
            <Route path="/" element={<MemeGallery />} />
            <Route path="/create" element={<CreateMeme />} />
            <Route path="/memes/:id" element={<MemeDetail />} />
            <Route path="/templates" element={<TemplateGallery />} />
            <Route path="/templates/create" element={<CreateTemplate />} />
          </Routes>
        </main>
      </div>
    </Router>
  );
}

function MemeGallery() {
  const [memes, setMemes] = useState([]);

  useEffect(() => {
    // Fetch memes from backend
    fetch('/api/memes')
      .then(response => response.json())
      .then(data => setMemes(data))
      .catch(error => console.error('Error fetching memes:', error));
  }, []);

  return (
    <div>
      <h2>Meme Gallery</h2>
      {memes.length === 0 ? (
        <p>No memes created yet. <Link to="/create">Create your first meme!</Link></p>
      ) : (
        <div className="meme-grid">
          {memes.map(meme => (
            <div key={meme.id} className="meme-item">
              <Link to={`/memes/${meme.id}`}>
                <img
                  src={`/memes/${meme.id}/image`}
                  alt={`${meme.template} meme`}
                  className="meme-thumbnail"
                  onError={(e) => {
                    // Fallback to a placeholder if image fails to load
                    e.target.src = 'data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iMjAwIiBoZWlnaHQ9IjIwMCIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj48cmVjdCB3aWR0aD0iMTAwJSIgaGVpZ2h0PSIxMDAlIiBmaWxsPSIjZGRkIi8+PHRleHQgeD0iNTAlIiB5PSI1MCUiIGZvbnQtZmFtaWx5PSJBcmlhbCIgZm9udC1zaXplPSIxOCIgZmlsbD0iIzk5OSIgdGV4dC1hbmNob3I9Im1pZGRsZSIgZHk9Ii4zZW0iPuKWiCAyMDB4MjAwPC90ZXh0Pjwvc3ZnPg==';
                  }}
                />
              </Link>
              <h3>{meme.template}</h3>
              <p>Top: {meme.text_top}</p>
              <p>Bottom: {meme.text_bottom}</p>
              <small>Created: {meme.created_at}</small>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}

function MemeDetail() {
  const { id } = useParams();
  const [meme, setMeme] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    // Fetch meme details from backend
    fetch(`/api/memes/${id}`)
      .then(response => {
        if (!response.ok) {
          throw new Error('Meme not found');
        }
        return response.json();
      })
      .then(data => {
        setMeme(data);
        setLoading(false);
      })
      .catch(error => {
        setError(error.message);
        setLoading(false);
      });
  }, [id]);

  if (loading) {
    return <div>Loading...</div>;
  }

  if (error) {
    return <div>Error: {error}</div>;
  }

  if (!meme) {
    return <div>Meme not found</div>;
  }

  return (
    <div className="meme-detail">
      <Link to="/">&larr; Back to Gallery</Link>
      <h2>Meme Details</h2>
      <div className="meme-full-image">
        <img
          src={`/memes/${meme.id}/image`}
          alt={`${meme.template} meme`}
          onError={(e) => {
            // Fallback to a placeholder if image fails to load
            e.target.src = 'data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iNjAwIiBoZWlnaHQ9IjQwMCIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj48cmVjdCB3aWR0aD0iMTAwJSIgaGVpZ2h0PSIxMDAlIiBmaWxsPSIjZGRkIi8+PHRleHQgeD0iNTAlIiB5PSI1MCUiIGZvbnQtZmFtaWx5PSJBcmlhbCIgZm9udC1zaXplPSIyNCIgZmlsbD0iIzk5OSIgdGV4dC1hbmNob3I9Im1pZGRsZSIgZHk9Ii4zZW0iPuKWiCBNZW1lIE5vdCBGb3VuZDwvdGV4dD48L3N2Zz4=';
          }}
        />
      </div>
      <div className="meme-info">
        <h3>{meme.template}</h3>
        <p><strong>Top Text:</strong> {meme.text_top}</p>
        <p><strong>Bottom Text:</strong> {meme.text_bottom}</p>
        <p><strong>Created:</strong> {meme.created_at}</p>
        <p><strong>ID:</strong> {meme.id}</p>
      </div>
    </div>
  );
}

function CreateMeme() {
  const [templates, setTemplates] = useState([]);
  const [selectedTemplate, setSelectedTemplate] = useState('');
  const [textTop, setTextTop] = useState('');
  const [textBottom, setTextBottom] = useState('');
  const [loading, setLoading] = useState(false);
  const [message, setMessage] = useState('');
  const navigate = useNavigate();

  useEffect(() => {
    // Fetch templates from backend
    fetch('/api/templates')
      .then(response => response.json())
      .then(data => setTemplates(data))
      .catch(error => console.error('Error fetching templates:', error));
  }, []);

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (!selectedTemplate) {
      setMessage('Please select a template');
      return;
    }
    
    setLoading(true);
    setMessage('');

    try {
      const response = await fetch('/api/memes', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          template: selectedTemplate,
          text_top: textTop,
          text_bottom: textBottom
        }),
      });

      if (response.ok) {
        const meme = await response.json();
        // Navigate to the meme detail page
        navigate(`/memes/${meme.id}`);
      } else {
        const error = await response.json();
        setMessage(`Error: ${error.error}`);
      }
    } catch (error) {
      setMessage(`Error: ${error.message}`);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="create-meme-container">
      <h2>Create New Meme</h2>
      {message && <div className="message">{message}</div>}
      <div className="create-meme-content">
        <div className="template-selector">
          <h3>Available Templates</h3>
          {templates.length === 0 ? (
            <p>No templates available. <Link to="/templates/create">Create a template</Link></p>
          ) : (
            <div className="template-list">
              {templates.map(template => (
                <div
                  key={template.name}
                  className={`template-option ${selectedTemplate === template.name ? 'selected' : ''}`}
                  onClick={() => setSelectedTemplate(template.name)}
                >
                  <img
                    src={`/templates/${template.name}/image`}
                    alt={`${template.name} template`}
                    className="template-preview"
                    onError={(e) => {
                      // Fallback to a placeholder if image fails to load
                      e.target.src = 'data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iMTAwIiBoZWlnaHQ9IjEwMCIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj48cmVjdCB3aWR0aD0iMTAwJSIgaGVpZ2h0PSIxMDAlIiBmaWxsPSIjZGRkIi8+PHRleHQgeD0iNTAlIiB5PSI1MCUiIGZvbnQtZmFtaWx5PSJBcmlhbCIgZm9udC1zaXplPSIxMiIgZmlsbD0iIzk5OSIgdGV4dC1hbmNob3I9Im1pZGRsZSIgZHk9Ii4zZW0iPlRlbXBsYXRlPC90ZXh0Pjwvc3ZnPg==';
                    }}
                  />
                  <span className="template-name">{template.name}</span>
                </div>
              ))}
            </div>
          )}
        </div>
        <div className="meme-form">
          <form onSubmit={handleSubmit}>
            <div>
              <label htmlFor="textTop">Top Text:</label>
              <input
                type="text"
                id="textTop"
                value={textTop}
                onChange={(e) => setTextTop(e.target.value)}
              />
            </div>
            <div>
              <label htmlFor="textBottom">Bottom Text:</label>
              <input
                type="text"
                id="textBottom"
                value={textBottom}
                onChange={(e) => setTextBottom(e.target.value)}
              />
            </div>
            <button type="submit" disabled={loading || !selectedTemplate}>
              {loading ? 'Creating...' : 'Create Meme'}
            </button>
          </form>
        </div>
      </div>
    </div>
  );
}

function TemplateGallery() {
  const [templates, setTemplates] = useState([]);

  useEffect(() => {
    // Fetch templates from backend
    fetch('/api/templates')
      .then(response => response.json())
      .then(data => setTemplates(data))
      .catch(error => console.error('Error fetching templates:', error));
  }, []);

  return (
    <div>
      <h2>Template Gallery</h2>
      {templates.length === 0 ? (
        <p>No templates created yet. <Link to="/templates/create">Create your first template!</Link></p>
      ) : (
        <div className="template-grid">
          {templates.map(template => (
            <div key={template.name} className="template-item">
              <Link to={`/templates/${template.name}`}>
                <img
                  src={`/templates/${template.name}/image`}
                  alt={`${template.name} template`}
                  className="template-thumbnail"
                  onError={(e) => {
                    // Fallback to a placeholder if image fails to load
                    e.target.src = 'data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iMjAwIiBoZWlnaHQ9IjIwMCIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj48cmVjdCB3aWR0aD0iMTAwJSIgaGVpZ2h0PSIxMDAlIiBmaWxsPSIjZGRkIi8+PHRleHQgeD0iNTAlIiB5PSI1MCUiIGZvbnQtZmFtaWx5PSJBcmlhbCIgZm9udC1zaXplPSIxOCIgZmlsbD0iIzk5OSIgdGV4dC1hbmNob3I9Im1pZGRsZSIgZHk9Ii4zZW0iPlRlbXBsYXRlPC90ZXh0Pjwvc3ZnPg==';
                  }}
                />
              </Link>
              <h3>{template.name}</h3>
              <small>Created: {template.created_at}</small>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}

function CreateTemplate() {
  const [name, setName] = useState('');
  const [image, setImage] = useState(null);
  const [loading, setLoading] = useState(false);
  const [message, setMessage] = useState('');

  const handleCreateTemplate = async (e) => {
    e.preventDefault();
    setLoading(true);
    setMessage('');

    try {
      // First create the template
      const response = await fetch('/api/templates', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ name }),
      });

      if (!response.ok) {
        const error = await response.json();
        throw new Error(error.error);
      }

      const template = await response.json();
      
      // If an image was selected, upload it
      if (image) {
        const formData = new FormData();
        formData.append('image', image);

        const imageResponse = await fetch(`/api/templates/${name}/image`, {
          method: 'POST',
          body: formData,
        });

        if (!imageResponse.ok) {
          const error = await imageResponse.json();
          throw new Error(error.error);
        }
      }

      setMessage(`Template "${template.name}" created successfully!`);
      // Reset form
      setName('');
      setImage(null);
    } catch (error) {
      setMessage(`Error: ${error.message}`);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div>
      <h2>Create New Template</h2>
      {message && <div className="message">{message}</div>}
      <form onSubmit={handleCreateTemplate}>
        <div>
          <label htmlFor="name">Template Name:</label>
          <input
            type="text"
            id="name"
            value={name}
            onChange={(e) => setName(e.target.value)}
            required
          />
        </div>
        <div>
          <label htmlFor="image">Template Image:</label>
          <input
            type="file"
            id="image"
            accept="image/*"
            onChange={(e) => setImage(e.target.files[0])}
          />
        </div>
        <button type="submit" disabled={loading}>
          {loading ? 'Creating...' : 'Create Template'}
        </button>
      </form>
    </div>
  );
}

export default App;