<?php


// ** MySQL settings ** //
/** The name of the database for WordPress */
define('DB_NAME', 'wordpress');

// Turn on debug functionalities 
define('WP_DEBUG', true);

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

define('AUTH_KEY',         'nDFJxT:s!?=A_?g&%1Z?ZeVlqGeEEF{UXcM^mYF~YKtFOmT=<biob}J{X])x.tNQ');
define('SECURE_AUTH_KEY',  '`83Rmf?xT-9ouQ3~iHlRte+t }!/^MdIJf6<0V]Qh,Q4#HGUjPl>~UC/%gpCUl;,');
define('LOGGED_IN_KEY',    'D-1+Nz2%,x`tHsNLST62(LzQOqH*{+%xam-do|6aPab qHCRg#I/`+u<B1cW|4bq');
define('NONCE_KEY',        '}s+]cha-k!NIvA3hgS<4-!]vl5cl=xr5+o,>>bC.I@hZU%Duw2} 33i^&w-P]/AJ');
define('AUTH_SALT',        'JSI+Vz`1zT,:(@;Xlum&F_HZwa%Te)@jkjA=F1t`W|lr29PTnmC+T2Rc7auQWnqN');
define('SECURE_AUTH_SALT', '#!+=x.pME[pjz@+Ymb_-p5:UR.}dh:A^do:n,aDK)f:F1Pa;-$B|FhXgpv5@S^=}');
define('LOGGED_IN_SALT',   'no9|)|dYUas0^)=VW@:yG#pKso-P}4pc>rz8H-%n+YjJDoL:E./bV>+[{r.[&],a');
define('NONCE_SALT',       ')N%1n8{5C_FGI3)u5&1Q54}({]5LC:0qAgdFOdCK#AwFKD&)Q^VL64{n|M2{x-s=');


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
