<?php

/**
 * The plugin bootstrap file
 *
 * This file is read by WordPress to generate the plugin information in the plugin
 * admin area. This file also includes all of the dependencies used by the plugin,
 * registers the activation and deactivation functions, and defines a function
 * that starts the plugin.
 *
 * @link              https://phz.fi/
 * @since             0.0.0
 * @package           Ric_Wp
 *
 * @wordpress-plugin
 * Plugin Name:       RIC WordPress Plugin
 * Plugin URI:        https://github.com/phzfi/RIC
 * Description:       This is the Responsive Image Cache integration plugin for WordPress. It replaces WP media URL's with RIC URL's.
 * Version:           0.0.1
 * Author:            Nicholas Saarela
 * Author URI:        https://phz.fi/
 * License:           GPL-2.0+
 * License URI:       http://www.gnu.org/licenses/gpl-2.0.txt
 * Text Domain:       ric-wp
 * Domain Path:       /languages
 */

// If this file is called directly, abort.
if ( ! defined( 'WPINC' ) ) {
	die;
}

/**
 * The code that runs during plugin activation.
 * This action is documented in includes/class-ric-wp-activator.php
 */
function activate_ric_wp() {
	require_once plugin_dir_path( __FILE__ ) . 'includes/class-ric-wp-activator.php';
	Ric_Wp_Activator::activate();
}

/**
 * The code that runs during plugin deactivation.
 * This action is documented in includes/class-ric-wp-deactivator.php
 */
function deactivate_ric_wp() {
	require_once plugin_dir_path( __FILE__ ) . 'includes/class-ric-wp-deactivator.php';
	Ric_Wp_Deactivator::deactivate();
}

register_activation_hook( __FILE__, 'activate_ric_wp' );
register_deactivation_hook( __FILE__, 'deactivate_ric_wp' );

/**
 * The core plugin class that is used to define internationalization,
 * admin-specific hooks, and public-facing site hooks.
 */
require plugin_dir_path( __FILE__ ) . 'includes/class-ric-wp.php';

/**
 * Begins execution of the plugin.
 *
 * Since everything within the plugin is registered via hooks,
 * then kicking off the plugin from this point in the file does
 * not affect the page life cycle.
 *
 * @since    1.0.0
 */
function run_ric_wp() {

	$plugin = new Ric_Wp();
	$plugin->run();

}

function load_js_file()
{
	wp_enqueue_script('client_js', plugins_url('/client.js',__FILE__));
}



class RicSettingsPage
{
    /**
     * Holds the values to be used in the fields callbacks
     */
    private $options;

    /**
     * Start up
     */
    public function __construct()
    {
        add_action( 'admin_menu', array( $this, 'add_plugin_page' ) );
        add_action( 'admin_init', array( $this, 'page_init' ) );
    }

    /**
     * Add options page
     */
    public function add_plugin_page()
    {
        // This page will be under "Settings"
        add_options_page(
            'Settings Admin', 
            'RIC Settings', 
            'manage_options', 
            'ric-setting-admin', 
            array( $this, 'create_admin_page' )
        );
    }

    /**
     * Options page callback
     */
    public function create_admin_page()
    {
        // Set class property
        $this->options = get_option( 'ric_option' );
        ?>
        <div class="wrap">
            <h2>RIC Wordpress Plugin Settings</h2>           
            <form method="post" action="options.php">
            <?php
                // This prints out all hidden setting fields
                settings_fields( 'ric_option_group' );   
                do_settings_sections( 'ric-setting-admin' );
                submit_button(); 
            ?>
            </form>
        </div>
        <?php
    }

    /**
     * Register and add settings
     */
    public function page_init()
    {        
        register_setting(
            'ric_option_group', // Option group
            'ric_option', // Option name
            array( $this, 'sanitize' ) // Sanitize
        );

        add_settings_section(
            'setting_section_id', // ID
            'RIC Server Settings', // Title
            array( $this, 'print_section_info' ), // Callback
            'ric-setting-admin' // Page
        );  

        add_settings_field(
            'ric_url', // ID
            'RIC Server URL', // Title 
            array( $this, 'id_number_callback' ), // Callback
            'ric-setting-admin', // Page
            'setting_section_id' // Section           
        );      

        /* add_settings_field(
            'title', 
            'Title', 
            array( $this, 'title_callback' ), 
            'ric-setting-admin', 
            'setting_section_id'
        ); 
        */    
    }

    /**
     * Sanitize each setting field as needed
     *
     * @param array $input Contains all settings fields as array keys
     */
    public function sanitize( $input )
    {
        $new_input = array();
        if( isset( $input['id_number'] ) )
            $new_input['id_number'] = absint( $input['id_number'] );

        if( isset( $input['title'] ) )
            $new_input['title'] = sanitize_text_field( $input['title'] );

        return $new_input;
    }

    /** 
     * Print the Section text
     */
    public function print_section_info()
    {
        print 'Enter your RIC Image Server URL below:';
    }

    /** 
     * Get the settings option array and print one of its values
     */
    public function id_number_callback()
    {
        printf(
            '<input type="text" id="id_number" name="ric_option[id_number]" value="%s" />',
            isset( $this->options['id_number'] ) ? esc_attr( $this->options['id_number']) : ''
        );
    }

    /** 
     * Get the settings option array and print one of its values
     */
   /* public function title_callback()
    {
        printf(
            '<input type="text" id="title" name="ric_option[title]" value="%s" />',
            isset( $this->options['title'] ) ? esc_attr( $this->options['title']) : ''
        );
    } */
}

if( is_admin() )
    $ric_settings_page = new RicSettingsPage();


/*

// SETTING UP THE SETTING MENUS 
/** Step 2 (from text above). 
add_action( 'admin_menu', 'ric_plugin_menu' );

/** Step 1. 
function ric_plugin_menu() {
	add_options_page( 'RIC Wordpress Plugin Options', 'RIC Wordpress Plugin', 'manage_options', 'ric-unique-identifier', 'ric_plugin_options' );
}

/** Step 3. 
function ric_plugin_options() {
	if ( !current_user_can( 'manage_options' ) )  {
		wp_die( __( 'You do not have sufficient permissions to access this page.' ) );
	}
	settings_fields( 'ricoption-group' );
	register_setting( 'ricoption-group', 'url-option');
	do_settings_sections( 'ricoption-group' );
	submit_button();
	<div class="wrap">
	<p>Here is where the form would gos if I actually had options.</p>';
	<form method="post" action="options.php">';
	do_settings_sections( 'ricoption-group' ); 
  <table class="form-table">';
  <tr valign="top">';
  <th scope="row">RIC Server URL</th>';
  <td><input type="text" name="url-option" value="  echo esc_attr( get_option('url-option') );
  echo '" /></td>';
  echo '</tr>';
	echo '</form>';
	echo '</div>';
}
*/


run_ric_wp();
add_action('wp_head', 'load_js_file');

?>