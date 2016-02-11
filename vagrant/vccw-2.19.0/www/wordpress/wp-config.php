<?php


// ** MySQL settings ** //
/** The name of the database for WordPress */
define('DB_NAME', 'wordpress');

/** MySQL database username */
define('DB_USER', 'wordpress');

/** MySQL database password */
define('DB_PASSWORD', 'wordpress');

/** MySQL hostname */
define('DB_HOST', 'localhost');

/** Database Charset to use in creating database tables. */
define('DB_CHARSET', 'utf8');

/** The Database Collate type. Don't change this if in doubt. */
define('DB_COLLATE', '');

define('AUTH_KEY',         '%g.x6^S:I|wyU?4YnL$sY2SqFIz-P-D(|=[7;`Gc%PB`?F72RI)IsG@g-1R=*Q-j');
define('SECURE_AUTH_KEY',  '?op+B_Hd|BI-=p;`<pPnPi60Rs:jdTO2lU$;k84N~jcP!&{c%|eX>0NGyq1<j,L]');
define('LOGGED_IN_KEY',    '(-Eq+C@iO7zr$dznrT@+B#%T|gd-_> 1p}<sQPF(h/=A>Wj}F4[S~.P/et3=Q2#[');
define('NONCE_KEY',        'VJya5hlDq,*|Pb,V>QhlA/Z;?2a#]^`;RR3i*.:8]DGP=Ld0B T5*j}c27/3)E]q');
define('AUTH_SALT',        'g,P1od>]sCzdBsM}.Z.?xFo:E2O;yuHWsAr]S*k=vLi`,D`q}7zd4W?g+m5B;qs ');
define('SECURE_AUTH_SALT', 'hUvXaB{= ,cy:Uz.|a*Lh?#;v^mL,(!6OeH`Gii?/BpRHrI~x:^9@!_ao&;53oDX');
define('LOGGED_IN_SALT',   '/bqlGv,)7 uel&e&r</N%`JcAE;x;)0k@&q6lVGp`>|2430|>^Fmk-|iL^ct/2r^');
define('NONCE_SALT',       '1<Zz GeKTqT>U]KWyw|c4j%:AbmneEf}Fb%8%%=[cuxkNK:?2n8eT ssQOfO{43|');


$table_prefix = 'wp_';


define( 'WP_HOME', 'http://vccw.dev' );
define( 'WP_SITEURL', 'http://vccw.dev' );
define( 'JETPACK_DEV_DEBUG', true );
define( 'WP_DEBUG', true );
define( 'FORCE_SSL_ADMIN', false );
define( 'SAVEQUERIES', false );



/* That's all, stop editing! Happy blogging. */

/** Absolute path to the WordPress directory. */
if ( !defined('ABSPATH') )
	define('ABSPATH', dirname(__FILE__) . '/');

/** Sets up WordPress vars and included files. */
require_once(ABSPATH . 'wp-settings.php');
