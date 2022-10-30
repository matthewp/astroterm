import '../styles/SiteNav.css';

function SiteLink({ children, link }) {
    return <li><a href={link}>{children}</a></li>
}

export default function() {
    return (
        <nav class="site-nav">
            <ul>
                <SiteLink link="https://github.com/matthewp/astroterm">GitHub</SiteLink>
            </ul>
        </nav>
    );
}